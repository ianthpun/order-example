package domain

type Asset interface {
}

type nftAsset struct {
}

func NewNFTAsset() nftAsset {
	return nftAsset{}
}
