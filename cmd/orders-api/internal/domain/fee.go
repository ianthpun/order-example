package domain

type Fee interface {
	GetAmount() string
	GetDescription() string
}

type fee struct {
	amount      string
	description string
}

var _ Fee = (*fee)(nil)

func NewFee(amount, description string) *fee {
	return &fee{
		amount:      amount,
		description: description,
	}
}

func NewNoFee() *fee {
	return &fee{
		amount:      "0",
		description: "no fee",
	}
}

func (f *fee) GetAmount() string {
	return f.amount
}

func (f *fee) GetDescription() string {
	return f.description
}
