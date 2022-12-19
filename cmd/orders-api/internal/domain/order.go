package domain

import (
	"fmt"
	"time"
)

type OrderState string

const (
	OrderStateCreated   OrderState = "CREATED"
	OrderStateExpired   OrderState = "EXPIRED"
	OrderStateCancelled OrderState = "CANCELLED"
	OrderStateConfirmed OrderState = "CONFIRMED"
)

func (o OrderState) String() string {
	return string(o)
}

type Order struct {
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

var orderExpiry = 10 * time.Minute

func NewOrder(
	ID string,
	userID string,
	asset Asset,
	price Money,
) (*Order, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be nil")
	}

	if asset == nil {
		return nil, fmt.Errorf("asset cannot be nil")
	}

	if price.IsZero() {
		return nil, fmt.Errorf("price cannot be zero")
	}

	return &Order{
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
) (Order, error) {
	var asset Asset
	var err error

	switch assetType {
	case AssetTypeDapperCredit.String():
		asset, err = NewDapperCreditAsset(NewMoney(amount, CurrencyType(currencyType)))
		if err != nil {
			return Order{}, err
		}
	case AssetTypeNFT.String():
		asset, err = NewNFTAsset(assetID, assetType)
		if err != nil {
			return Order{}, err
		}
	}

	o, err := NewOrder(
		ID,
		userID,
		asset,
		NewMoney(amount, CurrencyType(currencyType)),
	)
	if err != nil {
		return Order{}, err
	}

	o.state = OrderState(state)

	return *o, nil
}

func (o *Order) GetPaymentOptions() []PaymentOption {
	return o.paymentOptions
}

func (o *Order) getPaymentOptionCharge(paymentType PaymentMethodType) (amount Money, fees Money, err error) {
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

func (o *Order) IsExpired() bool {
	if o.expiresAt.After(time.Now()) {
		o.state = OrderStateExpired
		return true
	}

	return false
}

func (o *Order) GetOrderState() OrderState {
	return o.state
}

func (o *Order) GetID() string {
	//TODO implement me
	panic("implement me")
}

func (o *Order) GetUserID() string {
	//TODO implement me
	panic("implement me")
}

func (o *Order) GetAsset() Asset {
	//TODO implement me
	panic("implement me")
}

func (o *Order) GetCreatedAt() time.Time {
	//TODO implement me
	panic("implement me")
}

func (o *Order) GetExpiresAt() time.Time {
	//TODO implement me
	panic("implement me")
}

func (o *Order) GetPrice() Money {
	//TODO implement me
	panic("implement me")
}

func (o *Order) GetSelectedPaymentOption() []PaymentOption {
	//TODO implement me
	panic("implement me")
}

func (o *Order) SelectPaymentOption(option PaymentOption) bool {
	//TODO implement me
	panic("implement me")
}

func (o *Order) GetStateChanges() []Message {
	//TODO implement me
	panic("implement me")
}

func (o *Order) Cancel() error {
	if o.state != OrderStateCreated {
		return fmt.Errorf("cannot cancel state is it is not in created state: %s", o.state)
	}

	o.state = OrderStateCancelled

	return nil
}

func (o *Order) Confirm(option PaymentOption) error {
	//TODO implement me
	panic("implement me")
}

func (o *Order) Expire() error {
	//TODO implement me
	panic("implement me")
}
