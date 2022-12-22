package spanner

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/uptrace/bun"
	"order-sample/cmd/orders-api/internal/domain"
)

type repo struct {
	db        *bun.DB
	publisher message.Publisher
}

var _ domain.OrderRepository = (*repo)(nil)

func NewOrderRepository() *repo {
	return &repo{}
}

func (r *repo) InsertNewOrder(ctx context.Context, order domain.Order) error {
	o := orderToSpannerModel(order)

	_, err := r.db.NewInsert().Model(&o).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) GetOrder(ctx context.Context, orderID string) (domain.Order, error) {
	o := Order{ID: orderID}
	_, err := r.db.NewSelect().Model(&o).WherePK().Relation("paymentOptions").Exec(ctx)
	if err != nil {
		return domain.Order{}, err
	}

	return domain.UnmarshalOrderFromDatabase(
		o.ID,
		o.UserID,
		o.AssetID,
		o.AssetName,
		o.State,
		o.Amount,
		o.CurrencyCode,
	)
}

func (r *repo) UpdateOrder(
	ctx context.Context, orderID string,
	updateFn func(ctx context.Context, order *domain.Order) (*domain.Order, error),
) error {
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		o := Order{ID: orderID}
		_, err := tx.NewSelect().Model(&o).WherePK().For("UPDATE").Exec(ctx)
		if err != nil {
			return err
		}

		domainOrder, err := domain.UnmarshalOrderFromDatabase(
			o.ID,
			o.UserID,
			o.AssetID,
			o.AssetName,
			o.State,
			o.Amount,
			o.CurrencyCode,
		)

		order, err := updateFn(ctx, &domainOrder)
		if err != nil {
			return err
		}

		updateOrder := orderToSpannerModel(*order)

		_, err = tx.NewUpdate().Model(&updateOrder).Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	})

}

func orderToSpannerModel(o domain.Order) Order {
	return Order{
		ID:             o.GetID(),
		AssetID:        o.GetAsset().GetID(),
		AssetName:      o.GetAsset().GetName(),
		State:          o.GetOrderState().String(),
		UserID:         o.GetUserID(),
		CreatedAt:      o.GetCreatedAt(),
		PaymentOptions: paymentOptionsToSpannerModel(o.GetPaymentOptions()),
	}
}

func paymentOptionsToSpannerModel(options []domain.PaymentOption) []*PaymentOption {
	var paymentOptions []*PaymentOption
	for _, opt := range options {
		paymentOptions = append(paymentOptions, paymentOptionToSpannerModel(opt))
	}

	return paymentOptions
}

func paymentOptionToSpannerModel(p domain.PaymentOption) *PaymentOption {
	return &PaymentOption{
		ID:       p.GetID(),
		OrderID:  p.GetOrderID(),
		Fee:      p.GetFee().GetAmount(),
		Subtotal: p.GetSubtotal().GetAmount(),
		Currency: p.GetCurrency().String(),
	}
}
