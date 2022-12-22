package domain

type CurrencyType string

const (
	CurrencyTypeUSD CurrencyType = "USD"
)

func (c CurrencyType) String() string {
	return string(c)
}

type Money struct {
	amount       string
	currencyType CurrencyType
}

func NewMoney(amount string, currencyType CurrencyType) Money {
	return Money{
		amount:       amount,
		currencyType: currencyType,
	}
}

func (m Money) GetAmount() string {
	return m.amount
}

func (m Money) GetCurrencyType() CurrencyType {
	return m.currencyType
}

func (m Money) IsZero() bool {
	return m.amount == "0" || m.amount == ""
}
