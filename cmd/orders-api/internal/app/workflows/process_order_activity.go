package workflows

import (
	"context"
	"order-sample/cmd/orders-api/internal/domain"
)

// ProcessOrderActivities is an activity that is used by our Temporal Workflow. Do not use this struct for other
// reasons.
type ProcessOrderActivities struct {
	PaymentService  PaymentService
	AssetService    AssetService
	OrderRepository domain.OrderRepository
}

// NewProcessOrderActivities returns a temporal Activity used for Temporal workflows. Do not use this function for
// other reasons.
func NewProcessOrderActivities(
	paymentService PaymentService,
	assetService AssetService,
	orderRepository domain.OrderRepository,
) *ProcessOrderActivities {
	return &ProcessOrderActivities{
		PaymentService:  paymentService,
		AssetService:    assetService,
		OrderRepository: orderRepository,
	}
}

func (a *ProcessOrderActivities) ChargePayment(
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

func (a *ProcessOrderActivities) CancelOrder(ctx context.Context, orderID string) error {
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

func (a *ProcessOrderActivities) DeliverOrder(ctx context.Context, orderID string) error {
	order, err := a.OrderRepository.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}

	return a.AssetService.Deliver(ctx, order)
}

func (a *ProcessOrderActivities) ConfirmOrder(
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

func (a *ProcessOrderActivities) ExpireOrder(
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

func (a *ProcessOrderActivities) RefundPayment(
	ctx context.Context,
	paymentChargeID string,
) error {
	return nil
}

func (a *ProcessOrderActivities) CreateOrder(
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
