package app

import (
	"context"
	"order-sample/cmd/orders-api/internal/domain"
)

type CancelOrderHandler CommandHandler[string]

type cancelOrderUseCase struct {
	paymentService  PaymentService
	assetService    AssetService
	orderWorkflow   OrderWorkflow
	orderRepository domain.OrderRepository
}

func NewCancelOrderHandler(
	paymentService PaymentService,
	assetService AssetService,
	orderRepository domain.OrderRepository,
) *cancelOrderUseCase {
	return &cancelOrderUseCase{
		paymentService:  paymentService,
		assetService:    assetService,
		orderRepository: orderRepository,
	}
}

func (c *cancelOrderUseCase) Handle(ctx context.Context, orderID string) error {
	if err := c.orderWorkflow.ConfirmOrder(ctx, orderID); err != nil {
		return err
	}

	return nil
}
