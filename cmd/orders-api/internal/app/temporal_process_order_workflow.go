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

// ProcessOrderRequest holds the request body for processing an order in the workflow
type ProcessOrderRequest struct {
	OrderID string
	UserID  string
	Asset   ProcessOrderAsset
	Price   ProcessOrderPrice
}

type ProcessOrderAsset struct {
	ID   string
	Type domain.AssetType
	Name string
}

type ProcessOrderPrice struct {
	Amount       string
	CurrencyType domain.CurrencyType
}

// TemporalProcessOrderWorkflow is a function specifically used for Temporal workflows. Do not use this function for
// other reasons.
func TemporalProcessOrderWorkflow(ctx workflow.Context, req ProcessOrderRequest) (state string, err error) {
	var (
		processOrderActivity TemporalProcessOrderActivity
	)

	asset, err := domain.NewNFTAsset(req.Asset.ID, req.Asset.Name)
	if err != nil {
		return ProcessOrderStateFailed, err
	}

	order, err := domain.NewOrder(
		req.OrderID,
		req.UserID,
		asset,
		domain.NewMoney(req.Price.Amount, req.Price.CurrencyType),
	)
	if err != nil {
		return ProcessOrderStateFailed, err
	}

	workflow.GetLogger(ctx).Info("order in the beginning", req)

	err = workflow.ExecuteLocalActivity(
		WithDefaultLocalActivityOptions(ctx),
		processOrderActivity.CreateOrder,
		order,
	).Get(ctx, nil)
	if err != nil {
		return ProcessOrderStateFailed, fmt.Errorf("failed to create order: %w", err)
	}

	order, err = waitForOrderDecision(ctx, *order)
	if err != nil {
		return ProcessOrderStateFailed, fmt.Errorf("failed waiting for order decicision: %w", err)
	}

	// if the order decision is not to continue, just return
	if order.GetOrderState() != domain.OrderStateConfirmed {
		workflow.GetLogger(ctx).Info("state not confirmed, exiting", order.GetOrderState().String())
		return order.GetOrderState().String(), nil
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

func waitForOrderDecision(ctx workflow.Context, order domain.Order) (*domain.Order, error) {
	// wait for confirm/cancel/expiry signal
	confirmOrderChannel := workflow.GetSignalChannel(ctx, SignalChannels.CONFIRM_ORDER_CHANNEL)
	cancelOrderChannel := workflow.GetSignalChannel(ctx, SignalChannels.CANCEL_ORDER_CHANNEL)

	workflow.GetLogger(ctx).Info("waiting for order decision")

	var (
		processOrderActivity TemporalProcessOrderActivity
		o                    domain.Order
	)

	var signalErr error
	{
		selector := workflow.NewSelector(ctx)

		selector.AddReceive(confirmOrderChannel, func(c workflow.ReceiveChannel, _ bool) {
			workflow.GetLogger(ctx).Info("confirmed")
			var event ConfirmOrderRequest
			c.Receive(ctx, &event)

			workflow.GetLogger(ctx).Info("order before confirmation", o)
			err := workflow.ExecuteLocalActivity(
				WithDefaultLocalActivityOptions(ctx),
				processOrderActivity.ConfirmOrder,
				order.GetID(),
				event.PaymentOptionID,
			).Get(ctx, &o)
			if err != nil {
				signalErr = err

				return
			}

			workflow.GetLogger(ctx).Info("order after confirmation", o)

			return
		})

		selector.AddReceive(cancelOrderChannel, func(c workflow.ReceiveChannel, _ bool) {
			workflow.GetLogger(ctx).Info("cancel")
			err := workflow.ExecuteLocalActivity(
				WithDefaultLocalActivityOptions(ctx),
				processOrderActivity.CancelOrder,
				order.GetID(),
			).Get(ctx, &o)
			if err != nil {
				signalErr = err

				return
			}

			return
		})

		selector.AddFuture(workflow.NewTimer(ctx, orderExpiryTime), func(f workflow.Future) {
			workflow.GetLogger(ctx).Info("expire")
			err := workflow.ExecuteLocalActivity(
				WithDefaultLocalActivityOptions(ctx),
				processOrderActivity.ExpireOrder,
				order.GetID(),
			).Get(ctx, &o)
			if err != nil {
				signalErr = err

				return
			}

			return
		})

		selector.Select(ctx)
	}

	if signalErr != nil {
		return nil, signalErr
	}

	return &o, nil
}
