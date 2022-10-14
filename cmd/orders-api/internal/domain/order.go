package domain

import (
	"fmt"
	"time"
)

type Order interface {
	GetSupportedPaymentMethods() []PaymentInstrumentType
	SetPaymentOptions([]PaymentInstrument) error
	IsExpired() bool
	GetOrderState() OrderState
}

type OrderState string

const (
	OrderStateCreated OrderState = "CREATED"
	OrderStateExpired OrderState = "EXPIRED"
)

type order struct {
	id             string
	asset          Asset
	state          OrderState
	userID         string
	amount         Money
	paymentOptions []PaymentOption
	expiresAt      time.Time
}

// this validates an order is always going to implement the Order interface
var _ Order = (*order)(nil)

var orderExpiry = 10 * time.Minute

func NewOrder(
	ID string,
	userID string,
	asset Asset,
	amount Money,
) (*order, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be nil")
	}

	if asset == nil {
		return nil, fmt.Errorf("asset cannot be nil")
	}

	if amount.IsZero() {
		return nil, fmt.Errorf("amount cannot be zero")
	}

	return &order{
		id:        ID,
		userID:    userID,
		asset:     asset,
		amount:    amount,
		state:     OrderStateCreated,
		expiresAt: time.Now().Add(orderExpiry),
	}, nil
}

func (o *order) GetSupportedPaymentMethods() []PaymentInstrumentType {
	switch o.asset.GetAssetType() {
	case AssetTypeDapperCredit:
		return []PaymentInstrumentType{
			PaymentInstrumentTypeCreditCard,
			PaymentInstrumentTypeCoinbaseCrypto,
		}

	case AssetTypeNFT:
		return []PaymentInstrumentType{
			PaymentInstrumentTypeCreditCard,
			PaymentInstrumentTypeDapperCredit,
		}

	default:
		return []PaymentInstrumentType{}
	}
}

func (o *order) SetPaymentOptions(options []PaymentInstrument) error {
	for _, option := range options {
		amount, fee, err := o.getPaymentOptionCharge(option.GetPaymentInstrumentType())
		if err != nil {
			return err
		}
		po, err := NewPaymentOption(option.GetPaymentInstrumentType(), amount, fee)
		if err != nil {
			return err
		}

		o.paymentOptions = append(o.paymentOptions, po)
	}

	return nil
}

func (o *order) getPaymentOptionCharge(paymentType PaymentInstrumentType) (amount Money, fees Money, err error) {
	switch paymentType {
	case PaymentInstrumentTypeCoinbaseCrypto, PaymentInstrumentTypeDapperCredit:
		amount = NewMoney(o.amount.GetAmount(), CurrencyTypeUSD)
		fees = NewMoney("0", CurrencyTypeUSD)
	case PaymentInstrumentTypeCreditCard:
		amount = NewMoney(o.amount.GetAmount(), CurrencyTypeUSD)
		fees = NewMoney("5.00", CurrencyTypeUSD)
	default:
		err = fmt.Errorf("unsupported payment type: %s", paymentType)
	}

	return
}

func (o *order) IsExpired() bool {
	if o.expiresAt.After(time.Now()) {
		o.state = OrderStateExpired
		return true
	}

	return false
}

func (o *order) GetOrderState() OrderState {
	return o.state
}
