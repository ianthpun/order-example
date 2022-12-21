package app_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"order-sample/cmd/orders-api/internal/app"
	"order-sample/cmd/orders-api/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UnitTestSuite) Test_ProcessOrder_Success() {
	asset, _ := domain.NewNFTAsset(uuid.NewString(), "cool doodle")
	order, _ := domain.NewOrder(
		uuid.NewString(),
		uuid.NewString(),
		*asset,
		domain.NewMoney("10.00", domain.CurrencyTypeUSD),
	)

	request := app.ProcessOrderRequest{
		OrderID: order.GetID(),
		UserID:  order.GetUserID(),
		Asset: app.ProcessOrderAsset{
			ID:   order.GetAsset().GetID(),
			Type: order.GetAsset().GetAssetType(),
			Name: order.GetAsset().GetName(),
		},
		Price: app.ProcessOrderPrice{
			Amount:       order.GetPrice().GetAmount(),
			CurrencyType: order.GetPrice().GetCurrencyType(),
		},
	}

	var activity app.TemporalProcessOrderActivity

	s.env.OnActivity(activity.CreateOrder, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, order *domain.Order) error {
			return nil
		})

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(
			app.SignalChannels.CONFIRM_ORDER_CHANNEL,
			app.ConfirmOrderRequest{
				OrderID:         order.GetID(),
				PaymentOptionID: order.GetPaymentOptions()[0].GetID(),
			},
		)

	}, time.Millisecond)

	s.env.OnActivity(activity.ConfirmOrder, mock.Anything, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, id string, paymentOptionID string) (*domain.Order, error) {
			o := *order
			err := o.ConfirmPaymentOption(paymentOptionID)
			assert.NoError(s.T(), err)
			return &o, nil
		})

	s.env.ExecuteWorkflow(
		app.TemporalProcessOrderWorkflow,
		request,
	)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

//func (s *UnitTestSuite) Test_ProcessOrder_Expired() {
//	asset, _ := domain.NewNFTAsset(uuid.NewString(), "cool doodle")
//	order, _ := domain.NewOrder(
//		uuid.NewString(),
//		uuid.NewString(),
//		asset,
//		domain.NewMoney("10.00", domain.CurrencyTypeUSD),
//	)
//
//	var activity app.TemporalProcessOrderActivity
//	s.env.OnActivity(activity.CreateOrder, mock.Anything, mock.Anything).Return(
//		func(ctx context.Context, order domain.Order) error {
//			return nil
//		})
//
//	s.env.OnActivity(activity.ExpireOrder, mock.Anything, mock.Anything).Return(
//		func(ctx context.Context, id string) (domain.Order, error) {
//			return domain.Order{}, nil
//		})
//
//	s.env.ExecuteWorkflow(
//		app.TemporalProcessOrderWorkflow,
//		app.ProcessOrderRequest{Order: *order},
//	)
//
//	s.True(s.env.IsWorkflowCompleted())
//	s.NoError(s.env.GetWorkflowError())
//
//	//var result string
//	//s.env.GetWorkflowResult(&result)
//	//s.Equal(domain.OrderStateExpired, result)
//}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
