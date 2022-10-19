package domain

import "fmt"

type Asset interface {
	GetAssetType() AssetType
	GetID() string
	GetName() string
	IsEmpty() bool
}

type AssetType string

const (
	AssetTypeNFT          AssetType = "NFT"
	AssetTypeDapperCredit AssetType = "DAPPER_CREDIT"
)

func (a AssetType) String() string {
	return string(a)
}

type dapperCreditAsset struct {
	id     string
	name   string
	amount Money
}

func (d dapperCreditAsset) IsEmpty() bool {
	return d == dapperCreditAsset{}
}

const (
	dapperCreditID = "DAPPER_CREDIT"
)

var _ Asset = (*dapperCreditAsset)(nil)

func NewDapperCreditAsset(amount Money) (*dapperCreditAsset, error) {
	if amount.GetCurrencyType() != CurrencyTypeUSD {
		return nil, fmt.Errorf("dapper credit can only be purchasable as USD")
	}

	if amount.IsZero() {
		return nil, fmt.Errorf("price cannot be zero")
	}

	return &dapperCreditAsset{
		id:     dapperCreditID,
		name:   fmt.Sprintf("Dapper Credit: %s", amount.GetAmount()),
		amount: amount,
	}, nil
}

func (d *dapperCreditAsset) GetAssetType() AssetType {
	return AssetTypeDapperCredit
}

func (d *dapperCreditAsset) GetID() string {
	return d.id
}

func (d *dapperCreditAsset) GetName() string {
	return d.name
}

type nftAsset struct {
	id          string
	name        string
	ownerUserID string
}

var _ Asset = (*nftAsset)(nil)

func NewNFTAsset(id string, name string) (*nftAsset, error) {
	return &nftAsset{
		id:   id,
		name: name,
	}, nil
}

func (n nftAsset) GetAssetType() AssetType {
	return AssetTypeNFT
}

func (n nftAsset) GetID() string {
	return n.id
}

func (n nftAsset) GetName() string {
	return n.name
}

func (n nftAsset) IsEmpty() bool {
	return n == nftAsset{}
}

func (n nftAsset) GetOwnerUserID() string {
	return n.ownerUserID
}
