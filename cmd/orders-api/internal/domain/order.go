package domain

import (
	"fmt"
	"time"
)

type Order interface {
	GetID() string
	GetPaymentOptions() []PaymentOption
	GetSelectedPaymentOption() []PaymentOption
	SelectPaymentOption(option PaymentOption) bool
	GetUserID() string
	IsExpired() bool
	GetOrderState() OrderState
	GetAsset() Asset
	GetCreatedAt() time.Time
	GetExpiresAt() time.Time
	GetPrice() Money
	GetStateChanges() []Message
}

type OrderState string

const (
	OrderStateCreated OrderState = "CREATED"
	OrderStateExpired OrderState = "EXPIRED"
)

func (o OrderState) String() string {
	return string(o)
}

type order struct {
	id             string
	asset          Asset
	state          OrderState
	selectedOption PaymentOption
	userID         string
	merchantID     string
	price          Money
	paymentOptions []PaymentOption
	expiresAt      time.Time
	stateChanges   []Message
}

// this validates an order is always going to implement the Order interface
var _ Order = (*order)(nil)

var orderExpiry = 10 * time.Minute

func NewOrder(
	ID string,
	userID string,
	asset Asset,
	price Money,
) (*order, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be nil")
	}

	if asset == nil {
		return nil, fmt.Errorf("asset cannot be nil")
	}

	if price.IsZero() {
		return nil, fmt.Errorf("price cannot be zero")
	}

	return &order{
		id:             ID,
		userID:         userID,
		asset:          asset,
		price:          price,
		state:          OrderStateCreated,
		expiresAt:      time.Now().Add(orderExpiry),
		paymentOptions: paymentOptions(ID, price, asset),
	}, nil
}

func paymentOptions(orderID string, price Money, asset Asset) []PaymentOption {
	switch asset.GetAssetType() {
	case AssetTypeDapperCredit:
		return []PaymentOption{
			NewPaymentOption(orderID, PaymentMethodTypeCreditCard, price, NewNoFee()),
			NewPaymentOption(orderID, PaymentMethodTypeCoinbaseCrypto, price, NewNoFee()),
		}

	case AssetTypeNFT:
		return []PaymentOption{
			NewPaymentOption(orderID, PaymentMethodTypeCreditCard, price, NewNoFee()),
			NewPaymentOption(orderID, PaymentMethodTypeDapperCredit, price, NewNoFee()),
		}

	default:
		return []PaymentOption{}
	}
}

// UnmarshalOrderFromDatabase will return an order. Use this function strictly for unmarshalling from the database,
// as it may set the order in an incorrect state.
func UnmarshalOrderFromDatabase(
	ID string,
	userID string,
	assetID string,
	assetType string,
	state string,
	amount string,
	currencyType string,
) (*order, error) {
	var asset Asset
	var err error

	switch assetType {
	case AssetTypeDapperCredit.String():
		asset, err = NewDapperCreditAsset(NewMoney(amount, CurrencyType(currencyType)))
		if err != nil {
			return nil, err
		}
	case AssetTypeNFT.String():
		asset, err = NewNFTAsset(assetID, assetType)
		if err != nil {
			return nil, err
		}
	}

	o, err := NewOrder(
		ID,
		userID,
		asset,
		NewMoney(amount, CurrencyType(currencyType)),
	)
	if err != nil {
		return nil, err
	}

	o.state = OrderState(state)

	return o, nil
}

func (o *order) GetPaymentOptions() []PaymentOption {
	return o.paymentOptions
}

func (o *order) getPaymentOptionCharge(paymentType PaymentMethodType) (amount Money, fees Money, err error) {
	switch paymentType {
	case PaymentMethodTypeCoinbaseCrypto, PaymentMethodTypeDapperCredit:
		amount = NewMoney(o.price.GetAmount(), CurrencyTypeUSD)
		fees = NewMoney("0", CurrencyTypeUSD)
	case PaymentMethodTypeCreditCard:
		amount = NewMoney(o.price.GetAmount(), CurrencyTypeUSD)
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

func (o *order) GetID() string {
	//TODO implement me
	panic("implement me")
}

func (o *order) GetUserID() string {
	//TODO implement me
	panic("implement me")
}

func (o *order) GetAsset() Asset {
	//TODO implement me
	panic("implement me")
}

func (o *order) GetCreatedAt() time.Time {
	//TODO implement me
	panic("implement me")
}

func (o *order) GetExpiresAt() time.Time {
	//TODO implement me
	panic("implement me")
}

func (o *order) GetPrice() Money {
	//TODO implement me
	panic("implement me")
}

func (o *order) GetSelectedPaymentOption() []PaymentOption {
	//TODO implement me
	panic("implement me")
}

func (o *order) SelectPaymentOption(option PaymentOption) bool {
	//TODO implement me
	panic("implement me")
}

func (o *order) GetStateChanges() []Message {
	//TODO implement me
	panic("implement me")
}
