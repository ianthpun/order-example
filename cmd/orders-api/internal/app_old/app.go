package app_old

import (
	"context"
	temporalsdk "go.temporal.io/sdk/client"
	"order-sample/cmd/orders-api/internal/adapters/repository/spanner"
	"order-sample/cmd/orders-api/internal/adapters/services/grpc"
	"order-sample/cmd/orders-api/internal/adapters/services/temporal"
	temporal2 "order-sample/cmd/orders-api/internal/api/temporal"
	"order-sample/cmd/orders-api/internal/app_old/order"
)

// Application provides all Application capabilities
type Application struct {
	Order OrderHandler
}

// OrderHandler provides all OrderHandler capabilities
type OrderHandler struct {
	CreateOrder  order.CreateOrderHandler
	ConfirmOrder order.ConfirmOrderHandler
	CancelOrder  order.CancelOrderHandler
}

// CommandHandler
// These allow for all usecases under application to be private structs and without the need of multiple interfaces
type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (*R, error)
}

func New(ctx context.Context, temporalClient temporalsdk.Client) Application {
	paymentService := grpc.NewPaymentService()
	assetService := grpc.NewAssetService()
	orderRepository := spanner.NewOrderRepository()
	workflowService := temporal.NewWorkflowService(
		temporalClient,
		temporal.ProcessOrderConfig{
			Activities: temporal2.NewProcessOrderActivities(
				paymentService,
				assetService,
				orderRepository,
			),
			WorkflowFunc: temporal2.ProcessOrderWorkflow,
		},
	)

	return Application{
		Order: OrderHandler{
			CreateOrder: order.NewCreateOrderHandler(
				assetService,
				workflowService,
			),
			ConfirmOrder: order.NewConfirmOrderHandler(
				workflowService,
			),
			CancelOrder: order.NewCancelOrderHandler(
				workflowService,
			),
		},
	}
}
