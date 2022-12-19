package app

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/workflow"
	"order-sample/cmd/orders-api/internal/domain"
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

// TemporalProcessOrderActivity is an activity that is used by our Temporal Workflow. Do not use this struct for other
// reasons.
type TemporalProcessOrderActivity struct {
	PaymentService  PaymentService
	AssetService    AssetService
	OrderRepository domain.OrderRepository
	app             Application
}

func (a *TemporalProcessOrderActivity) ChargePayment(
	ctx context.Context,
	orderID string,
	userID string,
	option domain.PaymentOption,
) (string, error) {
	paymentChargeID, err := a.PaymentService.ChargePayment(
		ctx,
		orderID,
		userID,
		option,
	)
	if err != nil {
		return "", err
	}

	return paymentChargeID, nil
}

func (a *TemporalProcessOrderActivity) CancelOrder(ctx context.Context, orderID string) (domain.Order, error) {
	var o domain.Order
	err := a.OrderRepository.UpdateOrder(
		ctx,
		orderID,
		func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			o = *order
			if err := o.Cancel(); err != nil {
				return nil, err
			}

			return &o, nil
		})
	if err != nil {
		return domain.Order{}, err
	}

	return o, nil
}

func (a *TemporalProcessOrderActivity) ConfirmOrder(
	ctx context.Context,
	orderID string,
	option domain.PaymentOption,
) (domain.Order, error) {
	var o domain.Order
	if err := a.OrderRepository.UpdateOrder(
		ctx,
		orderID,
		func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			o = *order
			if err := o.Confirm(option); err != nil {
				return nil, err
			}

			return &o, nil
		}); err != nil {
		return domain.Order{}, err
	}

	return o, nil
}

func (a *TemporalProcessOrderActivity) ExpireOrder(
	ctx context.Context,
	orderID string,
) (domain.Order, error) {
	var o domain.Order
	err := a.OrderRepository.UpdateOrder(
		ctx,
		orderID,
		func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			o = *order
			if err := o.Expire(); err != nil {
				return nil, err
			}

			return &o, nil
		})
	if err != nil {
		return domain.Order{}, err
	}

	return o, nil
}

func (a *TemporalProcessOrderActivity) RefundPayment(
	ctx context.Context,
	paymentChargeID string,
) error {
	return nil
}

// NewTemporalProcessOrderActivity returns a temporal Activity used for Temporal workflows. Do not use this function for
// other reasons.
func NewTemporalProcessOrderActivity(
	paymentService PaymentService,
	assetService AssetService,
	orderRepository domain.OrderRepository,
) *TemporalProcessOrderActivity {
	return &TemporalProcessOrderActivity{
		PaymentService:  paymentService,
		AssetService:    assetService,
		OrderRepository: orderRepository,
	}
}

func (a *TemporalProcessOrderActivity) CreateOrder(
	ctx context.Context,
	order domain.Order,
) error {
	return a.OrderRepository.InsertNewOrder(ctx, order)
}

type ProcessOrderRequest struct {
	Order domain.Order
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
		nil,
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
				WithDefaultLocalActivityOptions(ctx), nil,
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
				nil,
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
				nil,
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
		nil,
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
				nil,
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
