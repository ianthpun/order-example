package domain

type PaymentOption interface {
	GetID() string
	GetOrderID() string
	GetPaymentType() PaymentMethodType
	GetFee() Fee
	GetSubtotal() Money
	GetTotal() Money
	GetCurrency() CurrencyType
}

type paymentOption struct {
	orderID      string
	paymentType  PaymentMethodType
	instrumentID string
	fee          Fee
	amount       Money
	currency     CurrencyType
}

var _ PaymentOption = (*paymentOption)(nil)

func NewPaymentOption(
	orderID string,
	paymentType PaymentMethodType,
	amount Money,
	fee Fee,
) *paymentOption {
	return &paymentOption{
		orderID:     orderID,
		paymentType: paymentType,
		amount:      amount,
		currency:    amount.GetCurrencyType(),
		fee:         fee,
	}
}

func (p paymentOption) GetID() string {
	//TODO implement me
	panic("implement me")
}

func (p paymentOption) GetCurrency() CurrencyType {
	return p.currency
}

func (p paymentOption) GetPaymentType() PaymentMethodType {
	//TODO implement me
	panic("implement me")
}

func (p paymentOption) GetFee() Fee {
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

func (p paymentOption) GetOrderID() string {
	return p.orderID
}
