package app

import (
	"context"
	"order-sample/cmd/orders-api/internal/domain"
)

type Application struct {
	PaymentService  PaymentService
	AssetService    AssetService
	OrderRepository domain.OrderRepository
}

// New returns a new Application
func New(
	paymentService PaymentService,
	assetService AssetService,
	orderRepository domain.OrderRepository,
) *Application {
	return &Application{
		PaymentService:  paymentService,
		AssetService:    assetService,
		OrderRepository: orderRepository,
	}
}

func (a *Application) ChargePayment(
	ctx context.Context,
	orderID string,
) (string, error) {
	order, err := a.OrderRepository.GetOrder(ctx, orderID)
	if err != nil {
		return "", err
	}

	paymentChargeID, err := a.PaymentService.ChargePayment(
		ctx,
		order.GetID(),
		order.GetUserID(),
		order.GetSelectedPaymentOption(),
	)
	if err != nil {
		return "", err
	}

	return paymentChargeID, nil
}

func (a *Application) CancelOrder(ctx context.Context, orderID string) error {
	return a.OrderRepository.UpdateOrder(
		ctx,
		orderID,
		func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			if err := order.Cancel(); err != nil {
				return nil, err
			}

			return order, nil
		})
}

func (a *Application) RequestDelivery(ctx context.Context, orderID string) error {
	order, err := a.OrderRepository.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}

	return a.AssetService.RequestDelivery(ctx, order)
}

func (a *Application) ConfirmOrder(
	ctx context.Context,
	orderID string,
	paymentOptionID string,
) error {
	return a.OrderRepository.UpdateOrder(
		ctx,
		orderID,
		func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			if err := order.ConfirmPaymentOption(paymentOptionID); err != nil {
				return nil, err
			}

			return order, nil
		})
}

func (a *Application) ExpireOrder(
	ctx context.Context,
	orderID string,
) error {
	return a.OrderRepository.UpdateOrder(
		ctx,
		orderID,
		func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			if err := order.Expire(); err != nil {
				return nil, err
			}

			return order, nil
		})
}

func (a *Application) RefundPayment(
	ctx context.Context,
	paymentChargeID string,
) error {
	return nil
}

func (a *Application) CreateOrder(
	ctx context.Context,
	order domain.Order,
) error {
	return a.OrderRepository.InsertNewOrder(ctx, order)
}

func (a *Application) OrderDelivered(
	ctx context.Context,
	orderID string,
) error {
	return a.OrderRepository.UpdateOrder(
		ctx,
		orderID,
		func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			if err := order.Delivered(); err != nil {
				return nil, err
			}

			return order, nil
		})
}

func (a *Application) OrderDeliveryFailed(
	ctx context.Context,
	orderID string,
	reason string,
) error {
	return a.OrderRepository.UpdateOrder(
		ctx,
		orderID,
		func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			if err := order.FailedDelivery(reason); err != nil {
				return nil, err
			}

			return order, nil
		})
}
