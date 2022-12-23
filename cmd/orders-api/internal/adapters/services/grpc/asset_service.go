package grpc

import (
	"context"
	"github.com/dapperlabs/dibs/v2/eventclient"
	"github.com/dapperlabs/dibs/v2/kafka"
	"order-sample/cmd/orders-api/internal/domain"
	pkgeventclient "order-sample/internal/eventclient"
	"order-sample/internal/protobuf/events"
)

const (
	dapperUserID = "Dapper"
)

type assetService struct {
	eventclient eventclient.Provider
}

func NewAssetService() *assetService {
	return &assetService{}
}

func (as *assetService) IsAvailable(
	ctx context.Context,
	asset domain.Asset,
) (bool, error) {
	return true, nil
}

func (as *assetService) RequestDelivery(ctx context.Context, order domain.Order) error {
	switch order.GetAsset().GetAssetType() {
	case domain.AssetTypeDapperCredit:
		return as.newCreditTransferRequest(ctx, order)
	case domain.AssetTypeNFT:
		return nil
	}

	return nil
}

func (as *assetService) newCreditTransferRequest(ctx context.Context, order domain.Order) error {
	topic := pkgeventclient.TopicCreditTransferRequested
	eventID := kafka.NewEventID(order.GetID(), topic)

	return as.eventclient.ProduceProto(ctx, pkgeventclient.TopicCreditTransferRequested, eventclient.ProtoMessage{
		ID:  eventID,
		Key: eventID,
		Proto: &events.CreditTransferRequested{
			Header: &events.Header{
				EventName: pkgeventclient.CreditTransferRequestedEventName,
				EventId:   eventID,
				RequestId: order.GetID(),
			},
			SenderUserId:                 dapperUserID,
			RecipientUserId:              order.GetUserID(),
			Amount:                       order.GetPrice().BaseUnits(),
			ReserveReceivedCredits:       false,
			ReleaseSpentReservation:      false,
			Type:                         events.CreditTransferType_OFF_CHAIN,
			Description:                  "credit user for order delivery",
			IsManual:                     true,
			IsVisibleToUser:              true,
			TransferCategory:             events.CreditTransferCategory_CREDIT_TRANSFER_ORDER_DELIVERY,
			ReceivedCreditsNonRedeemable: false,
		},
	})
}
