package grpc

import (
	"context"
	"order-sample/cmd/orders-api/internal/app"
	"order-sample/cmd/orders-api/internal/domain"
)

// WorkflowExecutor are all the capabilities of a workflow executor
type WorkflowExecutor interface {
	StartOrder(ctx context.Context, order domain.Order) error
	CancelOrder(ctx context.Context, orderID string) error
	ConfirmOrder(ctx context.Context, orderID string, paymentOptionID string) error
}

type OrderService struct {
	workflowEngine WorkflowExecutor
	app            app.Application
}

func NewOrderService() *OrderService {
	return &OrderService{}
}
