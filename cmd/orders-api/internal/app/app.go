package app

import (
	"context"
	temporalsdk "go.temporal.io/sdk/client"
	"order-sample/cmd/orders-api/internal/adapters/repository/spanner"
	"order-sample/cmd/orders-api/internal/adapters/services/grpc"
	"order-sample/cmd/orders-api/internal/adapters/services/temporal"
)

type Application struct {
	CreateOrder  CreateOrderHandler
	ConfirmOrder ConfirmOrderHandler
	CancelOrder  CancelOrderHandler
}

// CommandHandler
//These allow for all usecases under application to be private structs and without the need of multiple interfaces
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

	return Application{
		CreateOrder: NewCreateOrderHandler(
			paymentService,
			assetService,
			orderRepository,
		),
		ConfirmOrder: NewConfirmOrderHandler(
			paymentService,
			assetService,
			orderRepository,
			temporal.NewWorkflowService(
				temporalClient,
				temporal.ProcessOrderConfig{
					Activities: NewTemporalProcessOrderActivity(
						paymentService,
						assetService,
						orderRepository,
					),
					WorkflowFunc: TemporalProcessOrderWorkflow,
				}),
		),
		CancelOrder: NewCancelOrderHandler(
			paymentService,
			assetService,
			orderRepository,
		),
	}
}
