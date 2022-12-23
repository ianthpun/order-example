package eventloops

import (
	"context"
	"errors"
	"github.com/dapperlabs/dibs/v2/eventclient"
	"github.com/golang/protobuf/proto"
	"go.temporal.io/api/serviceerror"
	"order-sample/internal/eventloopbase"
	"order-sample/internal/protobuf/events"
	"order-sample/internal/protobuf/temporal"
)

func (w *Worker) handleCreditTransferComplete(ctx context.Context, eventData []byte) *eventclient.EventProcessResult {
	event := &events.CreditTransferComplete{}
	if err := proto.Unmarshal(eventData, event); err != nil {
		return eventloopbase.Fail("error in unmarshal: %s", err)
	}

	if event.GetTransferCategory() == events.CreditTransferCategory_CREDIT_TRANSFER_ORDER_DELIVERY {
		if err := w.signalDeliveryCompleteWorkflow(ctx, event); err != nil {
			return eventloopbase.Fail("error in unmarshal: %s", err)
		}
	}

	return nil
}

func (w *Worker) signalDeliveryCompleteWorkflow(ctx context.Context, event *events.CreditTransferComplete) error {
	signal := temporal.WorkflowOrderDeliveryCompleteSignal{
		OrderId: event.GetHeader().GetRequestId(),
		Status:  temporal.WorkflowOrderDeliveryCompleteSignal_DELIVERY_STATUS_SUCCEEDED,
	}

	if !event.GetSucceeded() {
		signal.Status = temporal.WorkflowOrderDeliveryCompleteSignal_DELIVERY_STATUS_FAILED
		signal.FailureReason = event.GetDetail()
	}

	err := w.workflowExecutor.SignalWorkflow(
		ctx,
		"somesome",
		"",
		temporal.WorkflowSignal_WORKFLOW_SIGNAL_ORDER_DELIVERY_COMPLETE.String(),
		&signal,
	)
	if err != nil {
		// ignore NotFound errors
		var nferr *serviceerror.NotFound
		if !errors.As(err, &nferr) {
			return err
		}
	}

	return err
}
