package workflows

import (
	"fmt"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/multierr"
	"order-sample/cmd/orders-api/internal/domain"
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
func ProcessOrderWorkflow(ctx workflow.Context, order Order) (state string, err error) {
	log := workflow.GetLogger(ctx)
	defer func() {
		if err != nil {
			log.Error("failed to process workflow", err)
		}
	}()

	log.Info("create new order", order)

	var (
		processOrderActivity ProcessOrderActivities
	)

	err = workflow.ExecuteLocalActivity(
		WithDefaultLocalActivityOptions(ctx),
		processOrderActivity.CreateOrder,
		order,
	).Get(ctx, nil)
	if err != nil {
		return ProcessOrderStateFailed, fmt.Errorf("failed to create order: %w", err)
	}

	log.Info("wait for order decision")

	decision, err := waitForOrderDecision(ctx, order.OrderID)
	if err != nil {
		return ProcessOrderStateFailed, fmt.Errorf("failed waiting for order decicision: %w", err)
	}

	// if the order decision is not to continue, just return
	if decision != OrderDecisionConfirmed {
		log.Info("state not confirmed, exiting", decision)
		return decision, nil
	}

	log.Info("process the payment")
	paymentChargeID, err := processPayment(ctx, order.OrderID)
	if err != nil {
		return ProcessOrderStateFailed, fmt.Errorf("failed to process payment: %w", err)
	}

	// refund payment if anything fails moving forward
	defer func() {
		if err != nil {
			refundErr := workflow.ExecuteLocalActivity(
				WithDefaultLocalActivityOptions(ctx),
				nil,
				processOrderActivity.RefundPayment,
				paymentChargeID,
			).Get(ctx, nil)

			err = multierr.Append(err, refundErr)
		}
	}()

	log.Info("deliver the order")
	err = deliverOrder(ctx, order)
	if err != nil {
		return ProcessOrderStateFailed, fmt.Errorf("failed to deliver order: %s", err)
	}

	log.Info("order succeeded")

	return ProcessOrderStateSucceeded, nil
}

// wait for confirm/cancel signal
func waitForOrderDecision(ctx workflow.Context, orderID string) (string, error) {
	confirmOrderChannel := workflow.GetSignalChannel(ctx, SignalChannels.CONFIRM_ORDER_CHANNEL)
	cancelOrderChannel := workflow.GetSignalChannel(ctx, SignalChannels.CANCEL_ORDER_CHANNEL)

	var (
		processOrderActivity ProcessOrderActivities
		orderDecision        string
		signalErr            error
	)
	{
		selector := workflow.NewSelector(ctx)

		selector.AddReceive(confirmOrderChannel, func(c workflow.ReceiveChannel, _ bool) {
			var event ConfirmOrderSignal
			c.Receive(ctx, &event)

			err := workflow.ExecuteLocalActivity(
				WithDefaultLocalActivityOptions(ctx),
				processOrderActivity.ConfirmOrder,
				orderID,
				event.PaymentOptionID,
			).Get(ctx, nil)
			if err != nil {
				signalErr = err

				return
			}
			orderDecision = OrderDecisionConfirmed

			return
		})

		selector.AddReceive(cancelOrderChannel, func(c workflow.ReceiveChannel, _ bool) {
			err := workflow.ExecuteLocalActivity(
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
		selector.AddFuture(workflow.NewTimer(ctx, orderExpiryTime), func(f workflow.Future) {
			err := workflow.ExecuteLocalActivity(
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

func processPayment(ctx workflow.Context, orderID string) (string, error) {
	var (
		paymentChargeID      string
		processOrderActivity ProcessOrderActivities
	)

	// attempt to charge the order
	err := workflow.ExecuteLocalActivity(
		WithDefaultLocalActivityOptions(ctx),
		processOrderActivity.ChargePayment,
		orderID,
	).Get(ctx, &paymentChargeID)
	if err != nil {
		return "", fmt.Errorf("failed to create order: %w", err)
	}

	return paymentChargeID, nil
}

func deliverOrder(ctx workflow.Context, order Order) error {
	var (
		processOrderActivity ProcessOrderActivities
	)

	err := workflow.ExecuteLocalActivity(
		WithDefaultLocalActivityOptions(ctx),
		processOrderActivity.DeliverOrder,
		order.OrderID,
	).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
