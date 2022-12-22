package domain

type Fee struct {
	amount      string
	description string
}

func NewFee(amount, description string) Fee {
	return Fee{
		amount:      amount,
		description: description,
	}
}

func NewNoFee() Fee {
	return Fee{
		amount:      "0",
		description: "no Fee",
	}
}

func (f *Fee) GetAmount() string {
	return f.amount
}

func (f *Fee) GetDescription() string {
	return f.description
}
