package domain

type PaymentInstrument interface {
	GetPaymentInstrumentType() PaymentInstrumentType
}

type paymentInstrument struct {
	paymentInstrumentType PaymentInstrumentType
}

type PaymentInstrumentType string

const (
	PaymentInstrumentTypeDapperCredit   PaymentInstrumentType = "DAPPER_CREDIT"
	PaymentInstrumentTypeCreditCard     PaymentInstrumentType = "CREDIT_CARD"
	PaymentInstrumentTypeCoinbaseCrypto PaymentInstrumentType = "COINBASE_CRYPTO"
)

func NewPaymentInstrument(paymentType PaymentInstrumentType) paymentInstrument {
	return paymentInstrument{
		paymentInstrumentType: paymentType,
	}
}

func (p *paymentInstrument) GetPaymentInstrumentType() PaymentInstrumentType {
	return p.paymentInstrumentType
}
