package app

import (
	"fmt"
	"go.temporal.io/sdk/workflow"
	"time"
)

const (
	orderExpiryTime = time.Minute * 5
)

// WithDefaultLocalActivityOptions returns the default local activity
func WithDefaultLocalActivityOptions(ctx workflow.Context) workflow.Context {
	return workflow.WithLocalActivityOptions(ctx, workflow.LocalActivityOptions{
		ScheduleToCloseTimeout: time.Hour * 2,
		StartToCloseTimeout:    time.Second * 20,
	})
}

// TemporalProcessOrderWorkflow is a function specifically used for Temporal workflows. Do not use this function for
// other reasons.
func TemporalProcessOrderWorkflow(ctx workflow.Context, req ProcessOrderRequest) (err error) {
	var (
		processOrderActivity TemporalProcessOrderActivity
	)

	order := req.Order

	err = workflow.ExecuteLocalActivity(
		WithDefaultLocalActivityOptions(ctx),
		processOrderActivity.CreateOrder,
		order,
	).Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// wait for confirm/cancel/expiry signal
	confirmOrderChannel := workflow.GetSignalChannel(ctx, SignalChannels.CONFIRM_ORDER_CHANNEL)
	cancelOrderChannel := workflow.GetSignalChannel(ctx, SignalChannels.CANCEL_ORDER_CHANNEL)

	var signalErr error
	var checkout bool
	{
		selector := workflow.NewSelector(ctx)

		selector.AddReceive(confirmOrderChannel, func(c workflow.ReceiveChannel, _ bool) {
			var event ConfirmOrderRequest
			c.Receive(ctx, &event)

			err = workflow.ExecuteLocalActivity(
				WithDefaultLocalActivityOptions(ctx),
				processOrderActivity.ConfirmOrder,
				event.OrderID,
				event.PaymentOption,
			).Get(ctx, &order)
			if err != nil {
				signalErr = err

				return
			}

			checkout = true
			// success state, continue
			return
		})

		selector.AddReceive(cancelOrderChannel, func(c workflow.ReceiveChannel, _ bool) {
			err = workflow.ExecuteLocalActivity(
				WithDefaultLocalActivityOptions(ctx),
				processOrderActivity.CancelOrder,
				order.GetID(),
			).Get(ctx, &order)
			if err != nil {
				signalErr = err

				return
			}

			return
		})

		selector.AddFuture(workflow.NewTimer(ctx, orderExpiryTime), func(f workflow.Future) {
			err = workflow.ExecuteLocalActivity(
				WithDefaultLocalActivityOptions(ctx),
				processOrderActivity.ExpireOrder,
				order.GetID(),
			).Get(ctx, &order)
			if err != nil {
				signalErr = err

				return
			}

			return
		})

		selector.Select(ctx)
	}

	if signalErr != nil {
		return signalErr
	}

	if !checkout {
		// exit as its done
		return nil
	}

	var paymentChargeID string

	// attempt to charge the order
	err = workflow.ExecuteLocalActivity(
		WithDefaultLocalActivityOptions(ctx),
		processOrderActivity.ChargePayment,
		order.GetID(),
		order.GetUserID(),
		order.GetSelectedPaymentOption(),
	).Get(ctx, &paymentChargeID)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	paymentProcessedChannel := workflow.GetSignalChannel(ctx, SignalChannels.CONFIRM_ORDER_CHANNEL)

	{
		selector := workflow.NewSelector(ctx)

		selector.AddReceive(paymentProcessedChannel, func(c workflow.ReceiveChannel, _ bool) {
			var event ConfirmOrderRequest
			c.Receive(ctx, &event)

			err = workflow.ExecuteLocalActivity(
				WithDefaultLocalActivityOptions(ctx),
				processOrderActivity.ConfirmOrder,
				event.OrderID,
				event.PaymentOption,
			).Get(ctx, &order)
			if err != nil {
				signalErr = err

				return
			}

			// success state, continue
			return
		})

		selector.Select(ctx)
	}

	if signalErr != nil {
		return signalErr
	}

	// refund payment if anything fails moving forward
	defer func() {
		if err != nil {
			err = workflow.ExecuteLocalActivity(
				WithDefaultLocalActivityOptions(ctx),
				nil,
				processOrderActivity.RefundPayment,
				paymentChargeID,
			).Get(ctx, &order)
			if err != nil {
				signalErr = err

				return
			}

		}
	}()

	return nil
}
