package domain

type PaymentOption struct {
	orderID      string
	paymentType  PaymentMethodType
	instrumentID string
	fee          Fee
	amount       Money
	currency     CurrencyType
}

func NewPaymentOption(
	orderID string,
	paymentType PaymentMethodType,
	amount Money,
	fee Fee,
) PaymentOption {
	return PaymentOption{
		orderID:     orderID,
		paymentType: paymentType,
		amount:      amount,
		currency:    amount.GetCurrencyType(),
		fee:         fee,
	}
}

func (p PaymentOption) GetID() string {
	return p.instrumentID
}

func (p PaymentOption) GetOrderID() string {
	return p.orderID
}

func (p PaymentOption) GetCurrency() CurrencyType {
	return p.currency
}

func (p PaymentOption) GetFee() Fee {
	return p.fee
}

func (p PaymentOption) GetSubtotal() Money {
	return p.amount
}
