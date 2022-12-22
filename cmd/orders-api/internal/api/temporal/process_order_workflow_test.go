package temporal

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"order-sample/cmd/orders-api/internal/app"
	"order-sample/cmd/orders-api/internal/domain"
	"order-sample/internal/protobuf/orders"
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

	var app app.Application

	s.env.OnActivity(app.CreateOrder, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, order Order) error {
			return nil
		})

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(
			orders.WorkflowSignal_WORKFLOW_SIGNAL_CONFIRM_ORDER.String(),
			orders.WorkflowConfirmOrderSignal{
				OrderId:         request.OrderID,
				PaymentOptionId: uuid.NewString(),
			},
		)

	}, time.Millisecond)

	s.env.OnActivity(app.ConfirmOrder, mock.Anything, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, id string, paymentOptionID string) error {
			return nil
		})

	s.env.OnActivity(app.ChargePayment, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, orderID string) (string, error) {
			return uuid.NewString(), nil
		})

	s.env.OnActivity(app.DeliverOrder, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, orderID string) error {
			return nil
		})

	s.env.ExecuteWorkflow(
		ProcessOrderWorkflow,
		request,
	)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var result string
	s.env.GetWorkflowResult(&result)
	s.Equal(ProcessOrderStateSucceeded, result)
}

func (s *UnitTestSuite) Test_ProcessOrder_Expired() {
	request := testOrderRequest()

	var app app.Application

	s.env.OnActivity(app.CreateOrder, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, order Order) error {
			return nil
		})

	// let the timer run until expiry happens

	s.env.OnActivity(app.ExpireOrder, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, orderID string) error {
			return nil
		})

	s.env.ExecuteWorkflow(
		ProcessOrderWorkflow,
		request,
	)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var result string
	s.env.GetWorkflowResult(&result)
	s.Equal(OrderDecisionExpired, result)
}

func (s *UnitTestSuite) Test_ProcessOrder_Cancelled() {
	request := testOrderRequest()

	var app app.Application

	s.env.OnActivity(app.CreateOrder, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, order Order) error {
			return nil
		})

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(
			orders.WorkflowSignal_WORKFLOW_SIGNAL_CANCEL_ORDER.String(),
			orders.WorkflowCancelOrderSignal{
				OrderId: request.OrderID,
			},
		)

	}, time.Millisecond)

	s.env.OnActivity(app.CancelOrder, mock.Anything, mock.Anything).Return(
		func(ctx context.Context, orderID string) error {
			return nil
		})

	s.env.ExecuteWorkflow(
		ProcessOrderWorkflow,
		request,
	)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var result string
	s.env.GetWorkflowResult(&result)
	s.Equal(OrderDecisionCancelled, result)
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func testOrderRequest() Order {
	asset, err := domain.NewNFTAsset(uuid.NewString(), "cool doodle")
	if err != nil {
		panic(err)
	}
	order, err := domain.NewOrder(
		uuid.NewString(),
		uuid.NewString(),
		*asset,
		domain.NewMoney("10.00", domain.CurrencyTypeUSD),
	)
	if err != nil {
		panic(err)
	}

	return Order{
		OrderID: order.GetID(),
		UserID:  order.GetUserID(),
		Asset: OrderAsset{
			ID:   order.GetAsset().GetID(),
			Type: order.GetAsset().GetAssetType(),
			Name: order.GetAsset().GetName(),
		},
		Price: OrderPrice{
			Amount:       order.GetPrice().GetAmount(),
			CurrencyType: order.GetPrice().GetCurrencyType(),
		},
	}
}
