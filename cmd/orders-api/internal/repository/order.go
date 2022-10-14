package repository

import (
	"context"
	"errors"
	"order-sample/cmd/orders-api/internal/domain"
)

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
)

type OrderRepository interface {
	InsertNewOrder(ctx context.Context, order domain.Order) error
}
