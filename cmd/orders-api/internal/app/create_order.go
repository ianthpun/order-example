package app

import (
	"context"
	"errors"
	"fmt"
	"order-sample/cmd/orders-api/internal/domain"
)

type CreateOrderHandler QueryHandler[CreateOrderRequest, domain.Order]

type createOrderUseCase struct {
	paymentService  PaymentService
	assetService    AssetService
	orderRepository domain.OrderRepository
}

func NewCreateOrderHandler(
	paymentService PaymentService,
	assetService AssetService,
	orderRepository domain.OrderRepository,
) *createOrderUseCase {
	return &createOrderUseCase{
		paymentService:  paymentService,
		assetService:    assetService,
		orderRepository: orderRepository,
	}
}

type CreateOrderRequest struct {
	IdempotencyKey string
	UserID         string
	Asset          domain.Asset
	Price          domain.Money
}

// Handle attempts to create a new order
func (c *createOrderUseCase) Handle(ctx context.Context, req CreateOrderRequest) (domain.Order, error) {
	assetAvailable, err := c.assetService.IsAvailable(ctx, req.Asset)
	if err != nil {
		return nil, fmt.Errorf("failed to check if asset was available: %s", err)
	}

	if !assetAvailable {
		return nil, fmt.Errorf("asset is not available for order")
	}

	order, err := domain.NewOrder(req.IdempotencyKey, req.UserID, req.Asset, req.Price)
	if err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	if err := c.orderRepository.InsertNewOrder(ctx, order); err != nil {
		if errors.Is(err, domain.ErrOrderAlreadyExists) {
			return nil, fmt.Errorf("order already exists")
		}

		return nil, err
	}

	return order, nil
}
