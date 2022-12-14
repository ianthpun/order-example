package domain

import "fmt"

const (
	AssetTypeNFT          AssetType = "NFT"
	AssetTypeDapperCredit AssetType = "DAPPER_CREDIT"

	dapperCreditID = "DAPPER_CREDIT"
)

type Asset struct {
	id        string
	name      string
	assetType AssetType
}

type AssetType string

func (a AssetType) String() string {
	return string(a)
}

func (a Asset) IsEmpty() bool {
	return a == Asset{}
}

func (a Asset) GetAssetType() AssetType {
	return a.assetType
}

func (a Asset) GetID() string {
	return a.id
}

func (a Asset) GetName() string {
	return a.name
}

func NewDapperCreditAsset(amount Money) (*Asset, error) {
	if amount.GetCurrencyType() != CurrencyTypeUSD {
		return nil, fmt.Errorf("dapper credit can only be purchasable as USD")
	}

	if amount.IsZero() {
		return nil, fmt.Errorf("price cannot be zero")
	}

	return &Asset{
		id:        dapperCreditID,
		name:      fmt.Sprintf("Dapper Credit: %s", amount.GetAmount()),
		assetType: AssetTypeDapperCredit,
	}, nil
}

func NewNFTAsset(id string, name string) (*Asset, error) {
	return &Asset{
		id:        id,
		name:      name,
		assetType: AssetTypeNFT,
	}, nil
}
