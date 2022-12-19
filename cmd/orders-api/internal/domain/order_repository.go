package domain

import (
	"context"
	"errors"
)

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
)

type OrderRepository interface {
	InsertNewOrder(ctx context.Context, order Order) error
	GetOrder(ctx context.Context, orderID string) (Order, error)
	UpdateOrder(
		ctx context.Context,
		orderID string,
		updateFn func(ctx context.Context, order *Order) error,
	) error
}
