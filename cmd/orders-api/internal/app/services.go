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
	) ([]domain.PaymentInstrument, error)
	ChargePayment(
		ctx context.Context,
		orderID string,
		userID string,
		paymentOption domain.PaymentOption,
	) (string, error)
	RefundPayment(
		ctx context.Context,
		paymentID string,
	) (string, error)
}

// AssetService are all the capabilities of the asset service
type AssetService interface {
	IsAvailable(ctx context.Context, asset domain.Asset) (bool, error)
	Deliver(ctx context.Context, order domain.Order) error
}
