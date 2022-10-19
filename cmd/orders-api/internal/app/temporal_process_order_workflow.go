package app

import (
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

type ProcessOrderRequest struct {
}

// TemporalProcessOrderWorkflow is a function specifically used for Temporal workflows. Do not use this function for
// other reasons.
func TemporalProcessOrderWorkflow(ctx workflow.Context, _ ProcessOrderRequest) error {
	return nil
}
