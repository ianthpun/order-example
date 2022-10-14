package spanner

import (
	"context"
	"order-sample/cmd/orders-api/internal/domain"
)

type repository struct {
}

func NewOrderRepository() *repository {
	return &repository{}
}

func (r repository) InsertNewOrder(ctx context.Context, order domain.Order) error {
	//TODO implement me
	panic("implement me")
}
