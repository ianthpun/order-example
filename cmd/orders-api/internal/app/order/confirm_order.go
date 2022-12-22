package order

import (
	"context"
	"order-sample/cmd/orders-api/internal/app"
)

type ConfirmOrderHandler app.CommandHandler[ConfirmOrderRequest]

type confirmOrderUseCase struct {
	workflowService OrderWorkflow
}

func NewConfirmOrderHandler(
	workflowService OrderWorkflow,
) *confirmOrderUseCase {
	c := confirmOrderUseCase{
		workflowService: workflowService,
	}

	return &c
}

type ConfirmOrderRequest struct {
	OrderID         string
	PaymentOptionID string
}

func (c *confirmOrderUseCase) Handle(ctx context.Context, req ConfirmOrderRequest) error {
	if err := c.workflowService.ConfirmOrder(ctx, req.OrderID, req.PaymentOptionID); err != nil {
		return err
	}

	return nil
}
