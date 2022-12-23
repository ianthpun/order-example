package temporal

import (
	"context"
	temporalsdk "go.temporal.io/sdk/client"
	"order-sample/cmd/orders-api/internal/domain"
)

type workflowExecutor struct {
	client temporalsdk.Client
}

func NewWorkflowExecutor(
	client temporalsdk.Client,
) *workflowExecutor {
	return &workflowExecutor{client: client}
}

func (w workflowExecutor) StartOrder(ctx context.Context, order domain.Order) error {
	//TODO implement me
	panic("implement me")
}

func (w workflowExecutor) CancelOrder(ctx context.Context, orderID string) error {
	//TODO implement me
	panic("implement me")
}

func (w workflowExecutor) ConfirmOrder(ctx context.Context, orderID string, paymentOptionID string) error {
	//TODO implement me
	panic("implement me")
}
