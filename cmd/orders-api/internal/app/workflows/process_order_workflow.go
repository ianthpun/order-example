package workflows

import (
	"fmt"
	workflowsdk "go.temporal.io/sdk/workflow"
	"go.uber.org/multierr"
	"order-sample/cmd/orders-api/internal/domain"
	"order-sample/internal/protobuf/orders"
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

// Order holds the request body for processing an order in the workflow
type Order struct {
	OrderID string
	UserID  string
	Asset   OrderAsset
	Price   OrderPrice
}

// OrderAsset holds the asset of an order
type OrderAsset struct {
	ID   string
	Type domain.AssetType
	Name string
}

// OrderPrice holds the price of an order
type OrderPrice struct {
	Amount       string
	CurrencyType domain.CurrencyType
}

// ProcessOrderWorkflow is a function specifically used for Temporal workflows. Do not use this function for
// other reasons.
func ProcessOrderWorkflow(ctx workflowsdk.Context, req *orders.WorkflowOrderRequest) (state string, err error) {
	log := workflowsdk.GetLogger(ctx)
	defer func() {
		if err != nil {
			log.Error("failed to process workflow", err)
		}
	}()

	log.Info("create new order", req)

	var (
		processOrderActivity ProcessOrderActivities
	)

	err = workflowsdk.ExecuteLocalActivity(
		WithDefaultLocalActivityOptions(ctx),
		processOrderActivity.CreateOrder,
		req,
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
				processOrderActivity.RefundPayment,
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
	confirmOrderChannel := workflowsdk.GetSignalChannel(ctx, SignalChannels.CONFIRM_ORDER_CHANNEL)
	cancelOrderChannel := workflowsdk.GetSignalChannel(ctx, SignalChannels.CANCEL_ORDER_CHANNEL)

	var (
		processOrderActivity ProcessOrderActivities
		orderDecision        string
		signalErr            error
	)
	{
		selector := workflowsdk.NewSelector(ctx)

		selector.AddReceive(confirmOrderChannel, func(c workflowsdk.ReceiveChannel, _ bool) {
			var event orders.WorkflowConfirmOrderSignal
			c.Receive(ctx, &event)

			err := workflowsdk.ExecuteLocalActivity(
				WithDefaultLocalActivityOptions(ctx),
				processOrderActivity.ConfirmOrder,
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
				processOrderActivity.CancelOrder,
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
				processOrderActivity.ExpireOrder,
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
		paymentChargeID      string
		processOrderActivity ProcessOrderActivities
	)

	// attempt to charge the order
	err := workflowsdk.ExecuteLocalActivity(
		WithDefaultLocalActivityOptions(ctx),
		processOrderActivity.ChargePayment,
		orderID,
	).Get(ctx, &paymentChargeID)
	if err != nil {
		return "", fmt.Errorf("failed to create order: %w", err)
	}

	return paymentChargeID, nil
}

func deliverOrder(ctx workflowsdk.Context, orderID string) error {
	var (
		processOrderActivity ProcessOrderActivities
	)

	err := workflowsdk.ExecuteLocalActivity(
		WithDefaultLocalActivityOptions(ctx),
		processOrderActivity.DeliverOrder,
		orderID,
	).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
