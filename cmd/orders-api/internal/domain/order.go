package domain

type Order interface {
}

type order struct {
}

func NewOrder() order {
	return order{}
}
