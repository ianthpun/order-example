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
}

// AssetService are all the capabilities of the asset service
type AssetService interface {
	IsAvailable(ctx context.Context, asset domain.Asset) (bool, error)
	Deliver(ctx context.Context, order domain.Order) error
}

// WorkflowExecutor are all the capabilities of a workflow executor
type WorkflowExecutor interface {
	StartOrder(ctx context.Context, order domain.Order)
	CancelOrder(ctx context.Context, orderID string)
	ConfirmOrder(ctx context.Context, orderID string, paymentOptionID string)
}
