package money

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
)

type Money struct {
	amount       decimal.Decimal
	currency     Currency
	withRounding *int32
}

type Currency string

const (
	CurrencyUSD  Currency = "USD"
	CurrencyFLOW Currency = "FLOW"
	CurrencyETH  Currency = "ETH"
	CurrencyBTC  Currency = "BTC"
)

func (c Currency) IsSupported() bool {
	return c == CurrencyFLOW || c == CurrencyETH || c == CurrencyBTC || c == CurrencyUSD
}

func defaultRounding(c Currency) *int32 {
	decimalPlaces := new(int32)

	switch c {
	case CurrencyUSD:
		*decimalPlaces = 2
	case CurrencyFLOW:
		*decimalPlaces = 8
	case CurrencyETH:
		*decimalPlaces = 3
	case CurrencyBTC:
		*decimalPlaces = 9
	default:
		return nil
	}

	return decimalPlaces
}

func NewFromString(amount string, curr Currency) (Money, error) {
	if !curr.IsSupported() {
		return Money{}, errors.New("currency not supported")
	}

	a, err := decimal.NewFromString(amount)
	if err != nil {
		return Money{}, err
	}

	return Money{
		amount:       a,
		currency:     curr,
		withRounding: defaultRounding(curr),
	}, nil
}

func (c *Money) Currency() Currency {
	return c.currency
}

func (c *Money) Amount() decimal.Decimal {
	return c.amount
}

func (c *Money) SetRounding(place int32) {
	*c.withRounding = place
}

func (c *Money) RemoveRounding() {
	c.withRounding = nil
}

func (c *Money) Add(b Money) error {
	if c.currency != b.currency {
		return errors.New("currency mismatch")
	}
	fmt.Println("before ", c.amount.String())

	fmt.Println("before b ", b.amount.String())

	c.amount = c.amount.Add(b.amount)
	if c.withRounding != nil {
		c.amount = c.amount.Round(*c.withRounding)
	} else {
		fmt.Println("nil rounding")
	}

	fmt.Println("after ", c.amount.String())

	return nil
}

func (c *Money) Sub(b Money) error {
	if c.currency != b.currency {
		return errors.New("currency mismatch")
	}

	c.amount = c.amount.Sub(b.amount)
	if c.withRounding != nil {
		c.amount = c.amount.Round(*c.withRounding)
	}

	return nil
}

func (c *Money) Div(b Money) error {
	if c.currency != b.currency {
		return errors.New("currency mismatch")
	}

	c.amount = c.amount.Div(b.amount)
	if c.withRounding != nil {
		c.amount = c.amount.Round(*c.withRounding)
	}

	return nil
}
