package grpc

import (
	"context"
	"order-sample/cmd/orders-api/internal/domain"
)

type assetService struct{}

func NewAssetService() *assetService {
	return &assetService{}
}

func (ps *assetService) IsAvailable(
	ctx context.Context,
	asset domain.Asset,
) (bool, error) {
	return true, nil
}
