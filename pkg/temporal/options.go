package temporal

import (
	"crypto/tls"

	"github.com/uber-go/tally/v4"
	"go.temporal.io/sdk/client"
	temporal_tally "go.temporal.io/sdk/contrib/tally"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"

	"go.temporal.io/server/common/log"
	"go.uber.org/zap"
)

// ClientOption is a function type use to set options for Temporal
type ClientOption func(s *ClientOptions)

// ClientOptions holds all available Temporal client options
type ClientOptions struct {
	client.Options
}

// WithHost sets the host:port of the Temporal server
func WithHost(host string) ClientOption {
	return func(s *ClientOptions) {
		s.HostPort = host
	}
}

// WithLogger sets a zap logger
func WithLogger(logger *zap.Logger) ClientOption {
	return func(s *ClientOptions) {
		s.Logger = log.NewSdkLogger(log.NewZapLogger(logger))
	}
}

// WithTLS sets a TLS configuration
func WithTLS(tls *tls.Config) ClientOption {
	return func(s *ClientOptions) {
		s.ConnectionOptions.TLS = tls
	}
}

// WithMetrics sets a metric scope
func WithMetrics(scope tally.Scope) ClientOption {
	return func(s *ClientOptions) {
		s.MetricsHandler = temporal_tally.NewMetricsHandler(scope)
	}
}

// WithNamespace sets the namespace of the Temporal server
func WithNamespace(namespace string) ClientOption {
	return func(s *ClientOptions) {
		s.Namespace = namespace
	}
}

// WorkerOption is a function type use to set options for Temporal workers
type WorkerOption func(s *WorkerOptions)

// WorkflowWithOptions provides a way to register a workflow with additional options
type WorkflowWithOptions struct {
	Workflow interface{}
	Options  workflow.RegisterOptions
}

// WorkerOptions holds all available Temporal Worker options
type WorkerOptions struct {
	worker.Options

	workflows  []interface{}
	activities []interface{}

	workflowsWithOptions []WorkflowWithOptions
}

// WithWorkflow adds a workflow that can be processed by a worker
func WithWorkflow(workflow interface{}) WorkerOption {
	return func(s *WorkerOptions) {
		s.workflows = append(s.workflows, workflow)
	}
}

// WithWorkflowWithOptions adds a workflow that can be processed by a worker.
// You could configure the workflow registration further with given option.
func WithWorkflowWithOptions(workflow interface{}, options workflow.RegisterOptions) WorkerOption {
	return func(s *WorkerOptions) {
		s.workflowsWithOptions = append(
			s.workflowsWithOptions,
			WorkflowWithOptions{
				Workflow: workflow,
				Options:  options,
			})
	}
}

// WithActivity adds an activity function/struct that can be processed by a worker
// If you registered a struct, it must be a pointer to the struct
func WithActivity(activity interface{}) WorkerOption {
	return func(s *WorkerOptions) {
		s.activities = append(s.activities, activity)
	}
}
