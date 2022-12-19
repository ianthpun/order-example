package app

import (
	"context"
	"order-sample/cmd/orders-api/internal/domain"
)

type ConfirmOrderHandler CommandHandler[ConfirmOrderRequest]

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
) *confirmOrderUseCase {
	c := confirmOrderUseCase{
		paymentService:  paymentService,
		assetService:    assetService,
		orderRepository: orderRepository,
	}

	return &c
}

type ConfirmOrderRequest struct {
	OrderID       string
	PaymentOption domain.PaymentOption
}

func (c *confirmOrderUseCase) Handle(ctx context.Context, req ConfirmOrderRequest) error {
	// TODO: run some validations first maybe

	if err := c.workflowService.ConfirmOrder(ctx, req.OrderID, req.PaymentOption); err != nil {
		return err
	}

	return nil
}
