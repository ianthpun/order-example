package grpc

import (
	"context"
	"order-sample/cmd/orders-api/internal/domain"
)

type paymentService struct{}

func (ps *paymentService) ChargePayment(ctx context.Context, orderID string, userID string, paymentOption domain.PaymentOption) (string, error) {
	//TODO implement me
	panic("implement me")
}

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
