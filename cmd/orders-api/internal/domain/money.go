package domain

type Money interface {
	GetAmount() string
	GetCurrencyType() CurrencyType
	IsZero() bool
}

type CurrencyType string

const (
	CurrencyTypeUSD  CurrencyType = "USD"
	CurrencyTypeFLOW CurrencyType = "FLOW"
)

func (c CurrencyType) String() string {
	return string(c)
}

var _ Money = (*money)(nil)

type money struct {
	amount       string
	currencyType CurrencyType
}

func NewMoney(amount string, currencyType CurrencyType) money {
	return money{
		amount:       amount,
		currencyType: currencyType,
	}
}

func (m money) GetAmount() string {
	return m.amount
}

func (m money) GetCurrencyType() CurrencyType {
	return m.currencyType
}

func (m money) IsZero() bool {
	//TODO implement me
	panic("implement me")
}
