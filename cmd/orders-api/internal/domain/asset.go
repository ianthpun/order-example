package domain

import "fmt"

type Asset interface {
	GetAssetType() AssetType
}

type AssetType string

const (
	AssetTypeNFT          AssetType = "NFT"
	AssetTypeDapperCredit AssetType = "DAPPER_CREDIT"
)

type dapperCreditAsset struct {
	amount Money
}

func NewDapperCreditAsset(amount Money) (*dapperCreditAsset, error) {
	if amount.GetCurrencyType() != CurrencyTypeUSD {
		return nil, fmt.Errorf("dapper credit can only be purchasable as USD")
	}

	if amount.IsZero() {
		return nil, fmt.Errorf("amount cannot be zero")
	}

	return &dapperCreditAsset{
		amount: amount,
	}, nil
}

func (d *dapperCreditAsset) GetAssetType() AssetType {
	return AssetTypeDapperCredit
}
