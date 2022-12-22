package workflows_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"order-sample/cmd/orders-api/internal/app/workflows"
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
	request := testOrderRequest()

	var activity workflows.ProcessOrderActivities

	s.env.OnActivity(activity.CreateOrder, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, order workflows.Order) error {
			return nil
		})

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(
			workflows.SignalChannels.CONFIRM_ORDER_CHANNEL,
			workflows.ConfirmOrderRequest{
				OrderID:         request.OrderID,
				PaymentOptionID: uuid.NewString(),
			},
		)

	}, time.Millisecond)

	s.env.OnActivity(activity.ConfirmOrder, mock.Anything, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, id string, paymentOptionID string) error {
			return nil
		})

	s.env.OnActivity(activity.ChargePayment, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, orderID string) (string, error) {
			return uuid.NewString(), nil
		})

	s.env.OnActivity(activity.DeliverOrder, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, orderID string) error {
			return nil
		})

	s.env.ExecuteWorkflow(
		workflows.ProcessOrderWorkflow,
		request,
	)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var result string
	s.env.GetWorkflowResult(&result)
	s.Equal(workflows.ProcessOrderStateSucceeded, result)
}

func (s *UnitTestSuite) Test_ProcessOrder_Expired() {
	request := testOrderRequest()

	var activity workflows.ProcessOrderActivities

	s.env.OnActivity(activity.CreateOrder, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, order workflows.Order) error {
			return nil
		})

	s.env.OnActivity(activity.ExpireOrder, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, orderID string) error {
			return nil
		})

	s.env.ExecuteWorkflow(
		workflows.ProcessOrderWorkflow,
		request,
	)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var result string
	s.env.GetWorkflowResult(&result)
	s.Equal(workflows.OrderDecisionExpired, result)
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func testOrderRequest() workflows.Order {
	asset, _ := domain.NewNFTAsset(uuid.NewString(), "cool doodle")
	order, _ := domain.NewOrder(
		uuid.NewString(),
		uuid.NewString(),
		*asset,
		domain.NewMoney("10.00", domain.CurrencyTypeUSD),
	)

	return workflows.Order{
		OrderID: order.GetID(),
		UserID:  order.GetUserID(),
		Asset: workflows.OrderAsset{
			ID:   order.GetAsset().GetID(),
			Type: order.GetAsset().GetAssetType(),
			Name: order.GetAsset().GetName(),
		},
		Price: workflows.OrderPrice{
			Amount:       order.GetPrice().GetAmount(),
			CurrencyType: order.GetPrice().GetCurrencyType(),
		},
	}
}