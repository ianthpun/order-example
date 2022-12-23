package grpc

import (
	"context"
	"order-sample/internal/protobuf/orders"
)

func (o *OrderService) CreateOrder(
	ctx context.Context,
	req *orders.CreateOrderRequest,
) (*orders.CreateOrderResponse, error) {
	return &orders.CreateOrderResponse{}, nil
}
