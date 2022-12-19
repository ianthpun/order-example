package app

import (
	"context"
	temporalsdk "go.temporal.io/sdk/client"
	"order-sample/cmd/orders-api/internal/adapters/repository/spanner"
	"order-sample/cmd/orders-api/internal/adapters/services/grpc"
	"order-sample/cmd/orders-api/internal/adapters/services/temporal"
)

type Application struct {
	Order OrderHandler
}

// OrderHandler provides all OrderHandler capabilities
type OrderHandler struct {
	CreateOrder  CreateOrderHandler
	ConfirmOrder ConfirmOrderHandler
	CancelOrder  CancelOrderHandler
}

// CommandHandler
// These allow for all usecases under application to be private structs and without the need of multiple interfaces
type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}

func New(ctx context.Context, temporalClient temporalsdk.Client) Application {
	paymentService := grpc.NewPaymentService()
	assetService := grpc.NewAssetService()
	orderRepository := spanner.NewOrderRepository()
	workflowService := temporal.NewWorkflowService(
		temporalClient,
		temporal.ProcessOrderConfig{
			Activities: NewTemporalProcessOrderActivity(
				paymentService,
				assetService,
				orderRepository,
			),
			WorkflowFunc: TemporalProcessOrderWorkflow,
		},
	)

	return Application{
		Order: OrderHandler{
			CreateOrder: NewCreateOrderHandler(
				paymentService,
				assetService,
				orderRepository,
				workflowService,
			),
			//ConfirmOrder: NewConfirmOrderHandler(
			//	paymentService,
			//	assetService,
			//	orderRepository,
			//),
			CancelOrder: NewCancelOrderHandler(
				paymentService,
				assetService,
				orderRepository,
			),
		},
	}
}
