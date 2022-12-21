package app

import (
	"fmt"
	"go.temporal.io/sdk/workflow"
	"order-sample/cmd/orders-api/internal/domain"
	"time"
)

const (
	orderExpiryTime = time.Minute * 5

	ProcessOrderStateFailed    = "FAILED"
	ProcessOrderStateSucceeded = "SUCCEEDED"
)

// WithDefaultLocalActivityOptions returns the default local activity
func WithDefaultLocalActivityOptions(ctx workflow.Context) workflow.Context {
	return workflow.WithLocalActivityOptions(ctx, workflow.LocalActivityOptions{
		ScheduleToCloseTimeout: time.Hour * 2,
		StartToCloseTimeout:    time.Second * 20,
	})
}

// Order holds the request body for processing an order in the workflow
type Order struct {
	OrderID string
	UserID  string
	Asset   OrderAsset
	Price   OrderPrice
}

type OrderAsset struct {
	ID   string
	Type domain.AssetType
	Name string
}

type OrderPrice struct {
	Amount       string
	CurrencyType domain.CurrencyType
}

// TemporalProcessOrderWorkflow is a function specifically used for Temporal workflows. Do not use this function for
// other reasons.
func TemporalProcessOrderWorkflow(ctx workflow.Context, order Order) (state string, err error) {
	var (
		processOrderActivity TemporalProcessOrderActivity
	)

	workflow.GetLogger(ctx).Info("order in the beginning", order)

	err = workflow.ExecuteLocalActivity(
		WithDefaultLocalActivityOptions(ctx),
		processOrderActivity.CreateOrder,
		order,
	).Get(ctx, nil)
	if err != nil {
		return ProcessOrderStateFailed, fmt.Errorf("failed to create order: %w", err)
	}

	decision, err := waitForOrderDecision(ctx, order.OrderID)
	if err != nil {
		return ProcessOrderStateFailed, fmt.Errorf("failed waiting for order decicision: %w", err)
	}

	// if the order decision is not to continue, just return
	if decision != domain.OrderStateConfirmed {
		workflow.GetLogger(ctx).Info("state not confirmed, exiting", decision)
		return decision.String(), nil
	}

	//
	//var paymentChargeID string
	//
	//// attempt to charge the order
	//err = workflow.ExecuteLocalActivity(
	//	WithDefaultLocalActivityOptions(ctx),
	//	processOrderActivity.ChargePayment,
	//	order.GetID(),
	//	order.GetUserID(),
	//	order.GetSelectedPaymentOption(),
	//).Get(ctx, &paymentChargeID)
	//if err != nil {
	//	return fmt.Errorf("failed to create order: %w", err)
	//}
	//
	//paymentProcessedChannel := workflow.GetSignalChannel(ctx, SignalChannels.CONFIRM_ORDER_CHANNEL)
	//
	//{
	//	selector := workflow.NewSelector(ctx)
	//
	//	selector.AddReceive(paymentProcessedChannel, func(c workflow.ReceiveChannel, _ bool) {
	//		var event ConfirmOrderRequest
	//		c.Receive(ctx, &event)
	//
	//		err = workflow.ExecuteLocalActivity(
	//			WithDefaultLocalActivityOptions(ctx),
	//			processOrderActivity.ConfirmOrder,
	//			event.OrderID,
	//			event.PaymentOption,
	//		).Get(ctx, &order)
	//		if err != nil {
	//			signalErr = err
	//
	//			return
	//		}
	//
	//		// success state, continue
	//		return
	//	})
	//
	//	selector.Select(ctx)
	//}
	//
	//if signalErr != nil {
	//	return signalErr
	//}
	//
	//// refund payment if anything fails moving forward
	//defer func() {
	//	if err != nil {
	//		err = workflow.ExecuteLocalActivity(
	//			WithDefaultLocalActivityOptions(ctx),
	//			nil,
	//			processOrderActivity.RefundPayment,
	//			paymentChargeID,
	//		).Get(ctx, &order)
	//		if err != nil {
	//			signalErr = err
	//
	//			return
	//		}
	//
	//	}
	//}()

	return ProcessOrderStateSucceeded, nil
}

func waitForOrderDecision(ctx workflow.Context, orderID string) (domain.OrderState, error) {
	// wait for confirm/cancel/expiry signal
	confirmOrderChannel := workflow.GetSignalChannel(ctx, SignalChannels.CONFIRM_ORDER_CHANNEL)
	cancelOrderChannel := workflow.GetSignalChannel(ctx, SignalChannels.CANCEL_ORDER_CHANNEL)

	workflow.GetLogger(ctx).Info("waiting for order decision")

	var (
		processOrderActivity TemporalProcessOrderActivity
		orderState           domain.OrderState
	)

	var signalErr error
	{
		selector := workflow.NewSelector(ctx)

		selector.AddReceive(confirmOrderChannel, func(c workflow.ReceiveChannel, _ bool) {
			workflow.GetLogger(ctx).Info("confirmed")
			var event ConfirmOrderRequest
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
			orderState = domain.OrderStateConfirmed

			return
		})

		selector.AddReceive(cancelOrderChannel, func(c workflow.ReceiveChannel, _ bool) {
			workflow.GetLogger(ctx).Info("cancel")
			err := workflow.ExecuteLocalActivity(
				WithDefaultLocalActivityOptions(ctx),
				processOrderActivity.CancelOrder,
				orderID,
			).Get(ctx, nil)
			if err != nil {
				signalErr = err

				return
			}

			orderState = domain.OrderStateCancelled

			return
		})

		selector.AddFuture(workflow.NewTimer(ctx, orderExpiryTime), func(f workflow.Future) {
			workflow.GetLogger(ctx).Info("expire")
			err := workflow.ExecuteLocalActivity(
				WithDefaultLocalActivityOptions(ctx),
				processOrderActivity.ExpireOrder,
				orderID,
			).Get(ctx, nil)
			if err != nil {
				signalErr = err

				return
			}

			orderState = domain.OrderStateExpired

			return
		})

		selector.Select(ctx)
	}

	if signalErr != nil {
		return "", signalErr
	}

	return orderState, nil
}
