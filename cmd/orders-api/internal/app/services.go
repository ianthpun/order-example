package app

import (
	"context"
	"order-sample/cmd/orders-api/internal/domain"
)

// PaymentService are all the capabilities of the payment service
type PaymentService interface {
	GetPaymentInstruments(
		ctx context.Context,
		userID string,
		types []domain.PaymentInstrumentType,
	) ([]domain.PaymentInstrument, error)
}

type AssetService interface {
	IsAvailable(ctx context.Context, asset domain.Asset) (bool, error)
}
