package temporal

import (
	"fmt"
	workflowsdk "go.temporal.io/sdk/workflow"
	"go.uber.org/multierr"
	"order-sample/cmd/orders-api/internal/domain"
	"order-sample/internal/protobuf/temporal"
	"time"
)

const (
	orderExpiryTime = time.Minute * 5

	ProcessOrderStateFailed    = "FAILED"
	ProcessOrderStateSucceeded = "SUCCEEDED"

	OrderDecisionCancelled = "CANCELLED"
	OrderDecisionConfirmed = "CONFIRMED"
	OrderDecisionExpired   = "EXPIRED"
)

// ProcessOrderWorkflow is a function specifically used for Temporal workflows. Do not use this function for
// other reasons.
func ProcessOrderWorkflow(ctx workflowsdk.Context, req *temporal.WorkflowOrderRequest) (state string, err error) {
	log := workflowsdk.GetLogger(ctx)
	defer func() {
		if err != nil {
			log.Error("failed to process workflow", err)
		}
	}()

	log.Info("create new order", req)
	createOrder, err := newOrder(req)
	if err != nil {
		if err != nil {
			return ProcessOrderStateFailed, fmt.Errorf("failed to generate new order: %w", err)
		}
	}

	var (
		app Activities
	)

	err = workflowsdk.ExecuteLocalActivity(
		WithDefaultLocalActivityOptions(ctx),
		app.CreateOrder,
		*createOrder,
	).Get(ctx, nil)
	if err != nil {
		return ProcessOrderStateFailed, fmt.Errorf("failed to create order: %w", err)
	}

	log.Info("wait for order decision")

	decision, err := waitForOrderDecision(ctx, req.GetOrderId())
	if err != nil {
		return ProcessOrderStateFailed, fmt.Errorf("failed waiting for order decision: %w", err)
	}

	// if the req decision is not to continue, just return
	if decision != OrderDecisionConfirmed {
		log.Info("state not confirmed, exiting", decision)
		return decision, nil
	}

	log.Info("process the payment")
	paymentChargeID, err := processPayment(ctx, req.GetOrderId())
	if err != nil {
		return ProcessOrderStateFailed, fmt.Errorf("failed to process payment: %w", err)
	}

	// refund payment if anything fails moving forward
	defer func() {
		if err != nil {
			refundErr := workflowsdk.ExecuteLocalActivity(
				WithDefaultLocalActivityOptions(ctx),
				nil,
				app.RefundPayment,
				paymentChargeID,
			).Get(ctx, nil)

			err = multierr.Append(err, refundErr)
		}
	}()

	log.Info("deliver the order")
	err = deliverOrder(ctx, req.GetOrderId())
	if err != nil {
		return ProcessOrderStateFailed, fmt.Errorf("failed to deliver order: %s", err)
	}

	log.Info("order succeeded")

	return ProcessOrderStateSucceeded, nil
}

// wait for confirm/cancel signal
func waitForOrderDecision(ctx workflowsdk.Context, orderID string) (string, error) {
	confirmOrderChannel := workflowsdk.GetSignalChannel(ctx, temporal.WorkflowSignal_WORKFLOW_SIGNAL_CONFIRM_ORDER.String())
	cancelOrderChannel := workflowsdk.GetSignalChannel(ctx, temporal.WorkflowSignal_WORKFLOW_SIGNAL_CANCEL_ORDER.String())

	var (
		app           Activities
		orderDecision string
		signalErr     error
	)
	{
		selector := workflowsdk.NewSelector(ctx)

		selector.AddReceive(confirmOrderChannel, func(c workflowsdk.ReceiveChannel, _ bool) {
			var event temporal.WorkflowConfirmOrderSignal
			c.Receive(ctx, &event)

			err := workflowsdk.ExecuteLocalActivity(
				WithDefaultLocalActivityOptions(ctx),
				app.ConfirmOrder,
				orderID,
				event.GetPaymentOptionId(),
			).Get(ctx, nil)
			if err != nil {
				signalErr = err

				return
			}
			orderDecision = OrderDecisionConfirmed

			return
		})

		selector.AddReceive(cancelOrderChannel, func(c workflowsdk.ReceiveChannel, _ bool) {
			err := workflowsdk.ExecuteLocalActivity(
				WithDefaultLocalActivityOptions(ctx),
				app.CancelOrder,
				orderID,
			).Get(ctx, nil)
			if err != nil {
				signalErr = err

				return
			}

			orderDecision = OrderDecisionCancelled

			return
		})
		// if no signal comes back before the order expiry time, expire the order
		selector.AddFuture(workflowsdk.NewTimer(ctx, orderExpiryTime), func(f workflowsdk.Future) {
			err := workflowsdk.ExecuteLocalActivity(
				WithDefaultLocalActivityOptions(ctx),
				app.ExpireOrder,
				orderID,
			).Get(ctx, nil)
			if err != nil {
				signalErr = err

				return
			}

			orderDecision = OrderDecisionExpired

			return
		})

		selector.Select(ctx)
	}

	if signalErr != nil {
		return "", signalErr
	}

	return orderDecision, nil
}

func processPayment(ctx workflowsdk.Context, orderID string) (string, error) {
	var (
		paymentChargeID string
		app             Activities
	)

	// attempt to charge the order
	err := workflowsdk.ExecuteLocalActivity(
		WithDefaultLocalActivityOptions(ctx),
		app.ChargePayment,
		orderID,
	).Get(ctx, &paymentChargeID)
	if err != nil {
		return "", fmt.Errorf("failed to create order: %w", err)
	}

	return paymentChargeID, nil
}

func deliverOrder(ctx workflowsdk.Context, orderID string) error {
	var (
		app Activities
	)

	err := workflowsdk.ExecuteLocalActivity(
		WithDefaultLocalActivityOptions(ctx),
		app.DeliverOrder,
		orderID,
	).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func newOrder(req *temporal.WorkflowOrderRequest) (*domain.Order, error) {
	asset, err := domain.NewDapperCreditAsset(
		domain.NewMoney(req.GetPrice().GetAmount(), domain.CurrencyTypeUSD),
	)
	if err != nil {
		return nil, err
	}

	return domain.NewOrder(
		req.GetOrderId(),
		req.GetUserId(),
		*asset,
		domain.NewMoney(req.GetPrice().GetAmount(), domain.CurrencyTypeUSD),
	)
}
