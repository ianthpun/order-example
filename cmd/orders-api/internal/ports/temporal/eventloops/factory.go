package eventloops

import (
	"github.com/dapperlabs/dibs/v2/eventclient"
	temporalsdk "go.temporal.io/sdk/client"
	"go.uber.org/zap"
	pkgeventclient "order-sample/internal/eventclient"
	"order-sample/internal/eventloopbase"
)

// Config defines the required configs for the orchestrator
type Config struct {
	Logger           *zap.Logger
	KafkaClient      eventclient.Provider
	EventClientConf  eventclient.Config
	Middleware       []eventloopbase.Middleware
	WorkflowExecutor temporalsdk.Client
}

type Worker struct {
	*eventloopbase.BaseEventLoop

	workflowExecutor temporalsdk.Client
}

// New creates a new payments event loop worker
func New(conf Config) Worker {
	w := Worker{
		BaseEventLoop:    eventloopbase.New(conf.Logger, conf.KafkaClient, conf.EventClientConf, conf.Middleware...),
		workflowExecutor: conf.WorkflowExecutor,
	}

	// Register event handlers
	w.Register(pkgeventclient.TopicCreditTransferComplete, w.handleCreditTransferComplete)

	return w
}
