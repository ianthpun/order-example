package app

import (
	"go.temporal.io/sdk/workflow"
	"order-sample/cmd/orders-api/internal/repository"
)

type TemporalProcessOrderActivity struct {
	paymentService  PaymentService
	assetService    AssetService
	orderRepository repository.OrderRepository
}

// NewTemporalProcessOrderActivity returns a temporal Activity used for Temporal workflows. Do not use this function for
// other means.
func NewTemporalProcessOrderActivity(
	paymentService PaymentService,
	assetService AssetService,
	orderRepository repository.OrderRepository,
) *TemporalProcessOrderActivity {
	return &TemporalProcessOrderActivity{
		paymentService:  paymentService,
		assetService:    assetService,
		orderRepository: orderRepository,
	}
}

// TemporalProcessOrderWorkflow is a function specifically used for Temporal workflows. Do not use this function for
// other means.
func TemporalProcessOrderWorkflow(ctx workflow.Context) error {
	return nil
}
