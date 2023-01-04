package temporal

import (
	"errors"
)

var (
	// ErrActivityNotImplemented is used as placeholder return for activity that hasn't implemented
	ErrActivityNotImplemented = errors.New("Activity is not implemented")
	// ErrWorkflowNotImplemented is used as placeholder return for workflow that hasn't implemented
	ErrWorkflowNotImplemented = errors.New("Workflow is not implemented")
	// ErrWorkflowShouldNotBeExecutedDirectly is used when a workflow should be executed from another workflow.
	ErrWorkflowShouldNotBeExecutedDirectly = errors.New("This workflow should not be executed directly")
)
