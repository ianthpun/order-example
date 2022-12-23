package temporal

import (
	temporalsdk "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type workflowService struct {
	worker worker.Worker
}

type ProcessOrderConfig struct {
	Activities   interface{}
	WorkflowFunc interface{}
}

func NewWorkflowService(
	client temporalsdk.Client,
	processOrder ProcessOrderConfig,
) *workflowService {
	// setup worker to be able to process available workflows
	w := worker.New(client, "orders", worker.Options{})

	w.RegisterActivity(processOrder.Activities)
	w.RegisterWorkflow(processOrder.WorkflowFunc)

	return &workflowService{worker: w}
}

func (w *workflowService) Run() error {
	return w.worker.Run(nil)
}
