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

func (a *TemporalProcessOrderActivity) ChargePayment(
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

func (a *TemporalProcessOrderActivity) CancelOrder(ctx context.Context, orderID string) error {
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

func (a *TemporalProcessOrderActivity) DeliverOrder(ctx context.Context, orderID string) error {
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

func (a *TemporalProcessOrderActivity) ConfirmOrder(
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

func (a *TemporalProcessOrderActivity) ExpireOrder(
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

func (a *TemporalProcessOrderActivity) RefundPayment(
	ctx context.Context,
	paymentChargeID string,
) error {
	return nil
}

func (a *TemporalProcessOrderActivity) CreateOrder(
	ctx context.Context,
	order Order,
) error {
	o, err := toOrderDomain(order)
	if err != nil {
		return err
	}

	return a.OrderRepository.InsertNewOrder(ctx, *o)
}

func toOrderDomain(o Order) (*domain.Order, error) {
	asset, err := domain.NewNFTAsset(o.Asset.ID, o.Asset.Name)
	if err != nil {
		return nil, err
	}

	return domain.NewOrder(
		o.OrderID,
		o.UserID,
		*asset,
		domain.NewMoney(o.Price.Amount, o.Price.CurrencyType),
	)
}

func (a *TemporalProcessOrderActivity) UpdateOrder(
	ctx context.Context,
	order domain.Order,
) (domain.Order, error) {
	var o domain.Order
	err := a.OrderRepository.UpdateOrder(
		ctx,
		order.GetID(),
		func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			return order, nil
		})
	if err != nil {
		return domain.Order{}, err
	}

	return o, nil
}
