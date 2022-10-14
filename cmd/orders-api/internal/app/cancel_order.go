package app

import (
	"context"
	"order-sample/cmd/orders-api/internal/repository"
)

type CancelOrderHandler CommandHandler[string]

type cancelOrderUseCase struct {
	paymentService  PaymentService
	assetService    AssetService
	orderRepository repository.OrderRepository
}

func NewCancelOrderHandler(
	paymentService PaymentService,
	assetService AssetService,
	orderRepository repository.OrderRepository,
) *cancelOrderUseCase {
	return &cancelOrderUseCase{
		paymentService:  paymentService,
		assetService:    assetService,
		orderRepository: orderRepository,
	}
}

func (c *cancelOrderUseCase) Handle(ctx context.Context, orderID string) error {
	return nil
}