package app

import (
	"context"
	"order-sample/cmd/orders-api/internal/domain"
)

// PaymentService are all the capabilities of the payment service
type PaymentService interface {
	GetPaymentMethods(
		ctx context.Context,
		userID string,
		types []domain.PaymentMethodType,
	) ([]domain.PaymentMethod, error)
	ChargePayment(
		ctx context.Context,
		orderID string,
		userID string,
		paymentOption domain.PaymentOption,
	) (string, error)
}

// AssetService are all the capabilities of the asset service
type AssetService interface {
	IsAvailable(ctx context.Context, asset domain.Asset) (bool, error)
}

// OrderWorkflow are all the capabilities of the order workflow
type OrderWorkflow interface {
	StartOrder(ctx context.Context, order domain.Order) error
	CancelOrder(ctx context.Context, orderID string) error
	ConfirmOrder(ctx context.Context, orderID string, paymentOption domain.PaymentOption) error
}
