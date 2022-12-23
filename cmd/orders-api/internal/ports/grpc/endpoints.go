package grpc

import (
	"context"
	"order-sample/cmd/orders-api/internal/domain"
	"order-sample/internal/protobuf/orders"
)

func (o *OrderService) CreateOrder(
	ctx context.Context,
	req *orders.CreateOrderRequest,
) (*orders.CreateOrderResponse, error) {
	order, err := newOrderRequest(req)
	if err != nil {
		return nil, err
	}

	if err := o.workflowExecutor.StartOrder(ctx, *order); err != nil {
		return nil, err
	}

	return &orders.CreateOrderResponse{}, nil
}

func newOrderRequest(req *orders.CreateOrderRequest) (*domain.Order, error) {
	asset, err := domain.NewDapperCreditAsset(domain.NewMoney(req.GetPrice().GetAmount(), domain.CurrencyTypeUSD))
	if err != nil {
		return nil, err
	}

	return domain.NewOrder(
		req.GetIdempotencyKey(),
		req.GetUserId(),
		*asset,
		domain.NewMoney(req.GetPrice().GetAmount(), domain.CurrencyTypeUSD),
	)
}

func (o *OrderService) CancelOrder(
	ctx context.Context,
	req *orders.CancelOrderRequest,
) (*orders.CancelOrderResponse, error) {
	if err := o.workflowExecutor.CancelOrder(ctx, req.GetOrderId()); err != nil {
		return nil, err
	}

	return &orders.CancelOrderResponse{}, nil
}

func (o *OrderService) ConfirmOrder(
	ctx context.Context,
	req *orders.ConfirmOrderRequest,
) (*orders.ConfirmOrderResponse, error) {
	if err := o.workflowExecutor.ConfirmOrder(ctx, req.GetOrderId(), req.GetPaymentOptionId()); err != nil {
		return nil, err
	}

	return &orders.ConfirmOrderResponse{}, nil
}
