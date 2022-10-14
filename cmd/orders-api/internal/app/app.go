package app

import (
	"order-sample/cmd/orders-api/internal/adapters/repository/spanner"
	"order-sample/cmd/orders-api/internal/adapters/services/grpc"
	"order-sample/cmd/orders-api/internal/repository"
)

type application struct {
	paymentService  PaymentService
	orderRepository repository.OrderRepository
}

func New() application {
	return application{
		paymentService:  grpc.NewPaymentService(),
		orderRepository: spanner.NewOrderRepository(),
	}
}
