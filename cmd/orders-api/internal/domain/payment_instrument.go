package domain

type PaymentMethod interface {
	GetID() string
	GetPaymentMethodType() PaymentMethodType
}

type paymentInstrument struct {
	id                    string
	paymentInstrumentType PaymentMethodType
}

var _ PaymentMethod = (*paymentInstrument)(nil)

type PaymentMethodType string

const (
	PaymentMethodTypeDapperCredit   PaymentMethodType = "DAPPER_CREDIT"
	PaymentMethodTypeCreditCard     PaymentMethodType = "CREDIT_CARD"
	PaymentMethodTypeCoinbaseCrypto PaymentMethodType = "COINBASE_CRYPTO"
)

func NewPaymentInstrument(id string, paymentType PaymentMethodType) paymentInstrument {
	return paymentInstrument{
		id:                    id,
		paymentInstrumentType: paymentType,
	}
}

func (p *paymentInstrument) GetPaymentMethodType() PaymentMethodType {
	return p.paymentInstrumentType
}

func (p *paymentInstrument) GetID() string {
	//TODO implement me
	panic("implement me")
}
