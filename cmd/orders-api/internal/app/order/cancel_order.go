package order

import (
	"context"
	"order-sample/cmd/orders-api/internal/app"
)

type CancelOrderHandler app.CommandHandler[string]

type cancelOrderUseCase struct {
	orderWorkflow OrderWorkflow
}

func NewCancelOrderHandler(
	orderWorkflow OrderWorkflow,
) *cancelOrderUseCase {
	return &cancelOrderUseCase{
		orderWorkflow: orderWorkflow,
	}
}

func (c *cancelOrderUseCase) Handle(ctx context.Context, orderID string) error {
	if err := c.orderWorkflow.CancelOrder(ctx, orderID); err != nil {
		return err
	}

	return nil
}
