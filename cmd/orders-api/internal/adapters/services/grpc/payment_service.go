package grpc

import (
	"context"
	"order-sample/cmd/orders-api/internal/domain"
)

type paymentService struct{}

func NewPaymentService() *paymentService {
	return &paymentService{}
}

func (ps *paymentService) GetPaymentInstruments(
	ctx context.Context,
	userID string,
	types []domain.PaymentInstrumentType,
) ([]domain.PaymentInstrument, error) {
	return []domain.PaymentInstrument{}, nil
}
