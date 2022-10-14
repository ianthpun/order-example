package domain

import "fmt"

type PaymentOption interface {
	GetPaymentType() PaymentInstrumentType
	GetFee() Money
	GetSubtotal() Money
	GetTotal() Money
}

type paymentOption struct {
	paymentType PaymentInstrumentType
	fees        Money
	amount      Money
}

func NewPaymentOption(
	paymentType PaymentInstrumentType,
	amount Money,
	fees Money,
) (*paymentOption, error) {
	if amount.IsZero() {
		return nil, fmt.Errorf("amount cannot be zero")
	}

	return &paymentOption{
		paymentType: paymentType,
		amount:      amount,
		fees:        fees,
	}, nil
}

func (p paymentOption) GetPaymentType() PaymentInstrumentType {
	//TODO implement me
	panic("implement me")
}

func (p paymentOption) GetFee() Money {
	//TODO implement me
	panic("implement me")
}

func (p paymentOption) GetSubtotal() Money {
	//TODO implement me
	panic("implement me")
}

func (p paymentOption) GetTotal() Money {
	//TODO implement me
	panic("implement me")
}
