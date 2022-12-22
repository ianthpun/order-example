package order

import (
	"context"
	"order-sample/cmd/orders-api/internal/domain"
)

// AssetService are all the capabilities of the asset service
type AssetService interface {
	IsAvailable(ctx context.Context, asset domain.Asset) (bool, error)
}

// OrderWorkflow are all the capabilities of the order workflow
type OrderWorkflow interface {
	StartOrder(ctx context.Context, order domain.Order) error
	CancelOrder(ctx context.Context, orderID string) error
	ConfirmOrder(ctx context.Context, orderID string, paymentOptionID string) error
}
