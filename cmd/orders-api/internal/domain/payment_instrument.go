package domain

type PaymentInstrument struct {
	id                    string
	paymentInstrumentType PaymentMethodType
}

type PaymentMethodType string

const (
	PaymentMethodTypeDapperCredit   PaymentMethodType = "DAPPER_CREDIT"
	PaymentMethodTypeCreditCard     PaymentMethodType = "CREDIT_CARD"
	PaymentMethodTypeCoinbaseCrypto PaymentMethodType = "COINBASE_CRYPTO"
)

func NewPaymentInstrument(id string, paymentType PaymentMethodType) PaymentInstrument {
	return PaymentInstrument{
		id:                    id,
		paymentInstrumentType: paymentType,
	}
}

func (p *PaymentInstrument) GetPaymentMethodType() PaymentMethodType {
	return p.paymentInstrumentType
}

func (p *PaymentInstrument) GetID() string {
	//TODO implement me
	panic("implement me")
}
