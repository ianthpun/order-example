package grpc

import (
	"context"
	"order-sample/cmd/orders-api/internal/app"
	"order-sample/cmd/orders-api/internal/domain"
)

// WorkflowExecutor are all the capabilities of a workflow executor
type WorkflowExecutor interface {
	StartOrder(ctx context.Context, order domain.Order)
	CancelOrder(ctx context.Context, orderID string)
	ConfirmOrder(ctx context.Context, orderID string, paymentOptionID string)
}

type OrderService struct {
	workflowExecutor WorkflowExecutor
	app              app.Application
}

func NewOrderService() *OrderService {
	return &OrderService{}
}
