package temporal

import (
	"fmt"
	temporalsdk "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type workflowService struct {
	client temporalsdk.Client
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

	go func() {
		if err := w.Run(nil); err != nil {
			fmt.Errorf("worker failed: %s", err)
		}
	}()

	return &workflowService{client: client}
}

func (w workflowService) RunProcessOrder() error {
	panic("implement me")
}
