package app

import (
	"context"
	"order-sample/cmd/orders-api/internal/domain"
)

type ConfirmOrderHandler CommandHandler[string]

type confirmOrderUseCase struct {
	paymentService  PaymentService
	assetService    AssetService
	orderRepository domain.OrderRepository
	workflowService OrderWorkflow
}

func NewConfirmOrderHandler(
	paymentService PaymentService,
	assetService AssetService,
	orderRepository domain.OrderRepository,
	workflowService OrderWorkflow,
) *confirmOrderUseCase {
	c := confirmOrderUseCase{
		paymentService:  paymentService,
		assetService:    assetService,
		orderRepository: orderRepository,
		workflowService: workflowService,
	}

	return &c
}

func (c *confirmOrderUseCase) Handle(ctx context.Context, id string) error {
	// TODO: run some validations first maybe

	if err := c.workflowService.ConfirmOrder(ctx, id); err != nil {
		return err
	}

	return nil
}
