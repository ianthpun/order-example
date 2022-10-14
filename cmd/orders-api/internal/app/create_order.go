package app

import (
	"context"
	"errors"
	"fmt"
	"order-sample/cmd/orders-api/internal/domain"
	"order-sample/cmd/orders-api/internal/repository"
)

type CreateOrderRequest struct {
	IdempotencyKey string
	UserID         string
	Asset          domain.Asset
	Amount         domain.Money
}

// CreateOrder attempts to create a new order
func (a application) CreateOrder(ctx context.Context, req CreateOrderRequest) (domain.Order, error) {
	assetAvailable, err := a.AssetService.IsAvailable(ctx, req.Asset)
	if err != nil {
		return nil, fmt.Errorf("failed to check if asset was available: %s", err)
	}

	if !assetAvailable {
		return nil, fmt.Errorf("asset is not available for order")
	}

	order, err := domain.NewOrder(req.IdempotencyKey, req.UserID, req.Asset, req.Amount)
	if err != nil {
		return nil, fmt.Errorf("error creating order: %w", err)
	}

	paymentInstruments, err := a.paymentService.GetPaymentInstruments(
		ctx,
		req.UserID,
		order.GetSupportedPaymentMethods(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment instruments from payment service")
	}

	if err := order.SetPaymentOptions(paymentInstruments); err != nil {
		return nil, err
	}

	if err := a.orderRepository.InsertNewOrder(ctx, order); err != nil {
		if errors.Is(err, repository.ErrOrderAlreadyExists) {
			return nil, fmt.Errorf("order already exists")
		}

		return nil, err
	}

	return order, nil
}
