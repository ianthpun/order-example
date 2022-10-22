package app

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/workflow"
	"order-sample/cmd/orders-api/internal/domain"
)

// TemporalProcessOrderActivity is an activity that is used by our Temporal Workflow. Do not use this struct for other
// reasons.
type TemporalProcessOrderActivity struct {
	PaymentService  PaymentService
	AssetService    AssetService
	OrderRepository domain.OrderRepository
}

// NewTemporalProcessOrderActivity returns a temporal Activity used for Temporal workflows. Do not use this function for
// other reasons.
func NewTemporalProcessOrderActivity(
	paymentService PaymentService,
	assetService AssetService,
	orderRepository domain.OrderRepository,
) *TemporalProcessOrderActivity {
	return &TemporalProcessOrderActivity{
		PaymentService:  paymentService,
		AssetService:    assetService,
		OrderRepository: orderRepository,
	}
}

func (a *TemporalProcessOrderActivity) CreateOrder(
	ctx context.Context,
	order domain.Order,
) error {
	return a.OrderRepository.InsertNewOrder(ctx, order)
}

type ProcessOrderRequest struct {
	Order domain.Order
}

// TemporalProcessOrderWorkflow is a function specifically used for Temporal workflows. Do not use this function for
// other reasons.
func TemporalProcessOrderWorkflow(ctx workflow.Context, req ProcessOrderRequest) error {
	var (
		processOrderActivity TemporalProcessOrderActivity
	)

	if err := workflow.ExecuteLocalActivity(
		//workflows.WithDefaultLocalActivityOptions(ctx),
		nil,
		processOrderActivity.CreateOrder,
		req.Order,
	).Get(ctx, nil); err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// wait for confirm/cancel/expiry signal

	return nil
}
