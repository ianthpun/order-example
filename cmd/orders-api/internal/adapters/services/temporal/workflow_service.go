package temporal

import (
	"context"
	"fmt"
	temporalsdk "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"order-sample/cmd/orders-api/internal/app"
	"order-sample/cmd/orders-api/internal/domain"
)

type workflowService struct {
	client temporalsdk.Client
}

var _ app.OrderWorkflow = (*workflowService)(nil)

type ProcessOrderConfig struct {
	Activities   interface{}
	WorkflowFunc interface{}
}

func NewWorkflowService(
	client temporalsdk.Client,
	processOrder ProcessOrderConfig,
) *workflowService {
	// setup worker to be able to process available workflows
	w := worker.New(client, "orders", worker.Options{})

	w.RegisterActivity(processOrder.Activities)
	w.RegisterWorkflow(processOrder.WorkflowFunc)

	go func() {
		if err := w.Run(nil); err != nil {
			fmt.Errorf("worker failed: %s", err)
		}
	}()

	return &workflowService{client: client}
}

func (w workflowService) StartOrder(ctx context.Context, order domain.Order) error {
	//TODO implement me
	panic("implement me")
}

func (w workflowService) CancelOrder(ctx context.Context, orderId string) error {
	//TODO implement me
	panic("implement me")
}

func (w workflowService) ConfirmOrder(ctx context.Context, orderID string) error {
	//TODO implement me
	panic("implement me")
}
