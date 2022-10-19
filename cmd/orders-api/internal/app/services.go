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
}

// AssetService are all the capabilities of the asset service
type AssetService interface {
	IsAvailable(ctx context.Context, asset domain.Asset) (bool, error)
}

// WorkflowService are all the capabilities of the workflow service
type WorkflowService interface {
	RunProcessOrder(ctx context.Context) error
}
