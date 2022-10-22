package money_test

import (
	"github.com/stretchr/testify/assert"
	"order-sample/pkg/money"
	"testing"
)

func TestMoney_Add(t *testing.T) {
	a, err := money.NewFromString("1.23", money.CurrencyUSD)
	assert.NoError(t, err)

	b, err := money.NewFromString("5.44", money.CurrencyUSD)
	assert.NoError(t, err)

	assert.NoError(t, a.Add(b))

	assert.Equal(t, "6.67", a.Amount())
}

func TestMoney_SetDecimalPlace(t *testing.T) {
	a, err := money.NewFromString("1.23", money.CurrencyUSD)
	assert.NoError(t, err)

	b, err := money.NewFromString("5.44", money.CurrencyUSD)
	assert.NoError(t, err)

	assert.NoError(t, a.Add(b))

	assert.Equal(t, "6.67", a.Amount().String())
}

func TestMoney_RemoveRounding(t *testing.T) {
	a, err := money.NewFromString("1.23", money.CurrencyUSD)
	assert.NoError(t, err)

	b, err := money.NewFromString("5.445", money.CurrencyUSD)
	assert.NoError(t, err)

	a.RemoveRounding()

	assert.NoError(t, a.Add(b))

	assert.Equal(t, "6.675", a.Amount())
}
