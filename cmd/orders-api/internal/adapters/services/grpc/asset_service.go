package grpc

import (
	"context"
	"order-sample/cmd/orders-api/internal/app"
	"order-sample/cmd/orders-api/internal/domain"
)

type assetService struct{}

var _ app.AssetService = (*assetService)(nil)

func NewAssetService() *assetService {
	return &assetService{}
}

func (ps *assetService) IsAvailable(
	ctx context.Context,
	asset domain.Asset,
) (bool, error) {
	return true, nil
}
