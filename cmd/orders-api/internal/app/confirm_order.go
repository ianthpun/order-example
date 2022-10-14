package app

import (
	"context"
	"order-sample/cmd/orders-api/internal/repository"
)

type ConfirmOrderHandler CommandHandler[string]

type confirmOrderUseCase struct {
	paymentService   PaymentService
	assetService     AssetService
	orderRepository  repository.OrderRepository
	workflowExecuter WorkflowService
}

func NewConfirmOrderHandler(
	paymentService PaymentService,
	assetService AssetService,
	orderRepository repository.OrderRepository,
	executer WorkflowService,
) *confirmOrderUseCase {
	c := confirmOrderUseCase{
		paymentService:   paymentService,
		assetService:     assetService,
		orderRepository:  orderRepository,
		workflowExecuter: executer,
	}

	return &c
}

func (c *confirmOrderUseCase) Handle(ctx context.Context, id string) error {
	if err := c.workflowExecuter.RunProcessOrder(); err != nil {
		return err
	}

	return nil
}
