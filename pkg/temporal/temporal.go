package temporal

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// NewClient returns a new Temporal client with options
func NewClient(opts ...ClientOption) (client.Client, error) {
	var options ClientOptions
	for _, o := range opts {
		o(&options)
	}

	c, err := client.NewLazyClient(options.Options)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// NewWorker returns a new Temporal worker with options
func NewWorker(client client.Client, taskQueue string, opts ...WorkerOption) worker.Worker {
	var options WorkerOptions
	for _, o := range opts {
		o(&options)
	}

	w := worker.New(client, taskQueue, options.Options)

	// since this could be a function name or a struct
	for i := range options.activities {
		w.RegisterActivity(options.activities[i])
	}

	for _, workflow := range options.workflows {
		w.RegisterWorkflow(workflow)
	}

	for _, workflowWithOption := range options.workflowsWithOptions {
		w.RegisterWorkflowWithOptions(
			workflowWithOption.Workflow,
			workflowWithOption.Options,
		)
	}

	return w
}
