package order

import (
	"context"
	"fmt"
	"order-sample/cmd/orders-api/internal/app"
	"order-sample/cmd/orders-api/internal/domain"
)

type CreateOrderHandler app.QueryHandler[CreateOrderRequest, domain.Order]

type createOrderUseCase struct {
	assetService  AssetService
	orderWorkflow OrderWorkflow
}

func NewCreateOrderHandler(
	assetService AssetService,
	orderWorkflow OrderWorkflow,
) *createOrderUseCase {
	return &createOrderUseCase{
		assetService:  assetService,
		orderWorkflow: orderWorkflow,
	}
}

type CreateOrderRequest struct {
	idempotencyKey string
	userID         string
	asset          domain.Asset
	price          domain.Money
}

func NewCreateOrderRequest(
	idempotencyKey string,
) CreateOrderRequest {
	return CreateOrderRequest{
		idempotencyKey: idempotencyKey,
	}
}

// Handle attempts to create a new order
func (c *createOrderUseCase) Handle(ctx context.Context, req CreateOrderRequest) (*domain.Order, error) {
	assetAvailable, err := c.assetService.IsAvailable(ctx, req.asset)
	if err != nil {
		return nil, fmt.Errorf("failed to check if asset was available: %s", err)
	}

	if !assetAvailable {
		return nil, fmt.Errorf("asset is not available for order")
	}

	order, err := domain.NewOrder(req.idempotencyKey, req.userID, req.asset, req.price)
	if err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	if err := c.orderWorkflow.StartOrder(ctx, *order); err != nil {
		return nil, err
	}

	return order, nil
}
