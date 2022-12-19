package app

import (
	"context"
	"order-sample/cmd/orders-api/internal/domain"
)

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
