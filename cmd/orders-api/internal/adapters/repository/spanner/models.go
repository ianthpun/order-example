package spanner

import "time"

// Order Represents the spanner model in the database
type Order struct {
	ID               string           `bun:"OrderId,pk"`
	AssetID          string           `bun:"AssetId"`
	AssetName        string           `bun:"AssetName"`
	AssetType        string           `bun:"AssetType"`
	State            string           `bun:"OrderState"`
	UserID           string           `bun:"UserId"`
	CreatedAt        time.Time        `bun:"CreatedAt"`
	UpdatedAt        time.Time        `bun:"UpdatedAt"`
	SelectedOptionId string           `bun:"SelectedOptionId"`
	PaymentOptions   []*PaymentOption `bun:"paymentOptions,rel:has-many,join:OrderId=OrderId"`
	Amount           string           `bun:"Amount"`
	CurrencyCode     string           `bun:"CurrencyCode"`
}

// PaymentOption represents the spanner model in the database
type PaymentOption struct {
	ID           string `bun:"PaymentOptionId"`
	OrderID      string `bun:"OrderId"`
	InstrumentID string `bun:"InstrumentId"`
	Subtotal     string `bun:"Subtotal"`
	Fee          string `bug:"Fee"`
	Currency     string `bun:"Currency"`
}
