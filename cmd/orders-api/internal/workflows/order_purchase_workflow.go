package workflows

type Executer interface {
	StartOrderWorkflow() error
}
