package domain

import (
	"fmt"
	"time"
)

type Order struct {
	id             string
	asset          Asset
	state          OrderState
	selectedOption PaymentOption
	userID         string
	merchantID     string
	price          Money
	paymentOptions []PaymentOption
	failureReason  string
	createdAt      time.Time
}

type OrderState string

const (
	OrderStateCreated        OrderState = "CREATED"
	OrderStateExpired        OrderState = "EXPIRED"
	OrderStateDelivered      OrderState = "DELIVERED"
	OrderStateDeliveryFailed OrderState = "DELIVERY_FAILED"
	OrderStateCancelled      OrderState = "CANCELLED"
	OrderStateConfirmed      OrderState = "CONFIRMED"
)

func (o OrderState) String() string {
	return string(o)
}

func NewOrder(
	ID string,
	userID string,
	asset Asset,
	price Money,
) (*Order, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be nil")
	}

	if asset.IsEmpty() {
		return nil, fmt.Errorf("asset cannot be empty")
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
		paymentOptions: paymentOptions(ID, price, asset),
		createdAt:      time.Now(),
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

func (o *Order) GetPaymentOptions() []PaymentOption {
	return o.paymentOptions
}

func (o *Order) GetOrderState() OrderState {
	return o.state
}

func (o *Order) GetID() string {
	return o.id
}

func (o *Order) GetUserID() string {
	return o.userID
}

func (o *Order) GetAsset() Asset {
	return o.asset
}

func (o *Order) GetCreatedAt() time.Time {
	return o.createdAt
}

func (o *Order) GetPrice() Money {
	return o.price
}

func (o *Order) GetSelectedPaymentOption() PaymentOption {
	return o.selectedOption
}

func (o *Order) Cancel() error {
	if o.state != OrderStateCreated {
		return fmt.Errorf("cannot cancel state is it is not in created state: %s", o.state)
	}

	o.state = OrderStateCancelled

	return nil
}

func (o *Order) ConfirmPaymentOption(paymentOptionID string) error {
	for _, po := range o.paymentOptions {
		if paymentOptionID == po.GetID() {
			o.selectedOption = po
			o.state = OrderStateConfirmed

			return nil
		}
	}

	return fmt.Errorf("paymentOptionID %s was not part of order", paymentOptionID)
}

func (o *Order) Expire() error {
	o.state = OrderStateExpired

	return nil
}

func (o *Order) Delivered() error {
	o.state = OrderStateDelivered

	return nil
}

func (o *Order) FailedDelivery(reason string) error {
	o.state = OrderStateDelivered
	o.failureReason = reason

	return nil
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
	var asset *Asset
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
		*asset,
		NewMoney(amount, CurrencyType(currencyType)),
	)
	if err != nil {
		return Order{}, err
	}

	o.state = OrderState(state)

	return *o, nil
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
