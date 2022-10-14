package domain

type PaymentInstrument interface {
	GetPaymentInstrumentType() PaymentInstrumentType
}

type paymentInstrument struct {
	paymentInstrumentType PaymentInstrumentType
}

type PaymentInstrumentType string

const (
	PaymentInstrumentTypeDapperCredit PaymentInstrumentType = "DAPPER_CREDIT"
)

func NewPaymentInstrument(paymentType PaymentInstrumentType) paymentInstrument {
	return paymentInstrument{
		paymentInstrumentType: paymentType,
	}
}

func (p *paymentInstrument) GetPaymentInstrumentType() PaymentInstrumentType {
	return p.paymentInstrumentType
}
