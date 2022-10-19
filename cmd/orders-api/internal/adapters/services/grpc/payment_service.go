package grpc

import (
	"context"
	"order-sample/cmd/orders-api/internal/app"
	"order-sample/cmd/orders-api/internal/domain"
)

type paymentService struct{}

var _ app.PaymentService = (*paymentService)(nil)

func NewPaymentService() *paymentService {
	return &paymentService{}
}

func (ps *paymentService) GetPaymentMethods(
	ctx context.Context,
	userID string,
	types []domain.PaymentMethodType,
) ([]domain.PaymentMethod, error) {
	return []domain.PaymentMethod{}, nil
}
