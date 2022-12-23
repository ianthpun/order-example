package middleware

import (
	"context"

	"github.com/dapperlabs/dibs/v3/eventclient"
	"github.com/golang/protobuf/proto"

	"github.com/dapperlabs/dapper-flow-api/pkg/protobuf/dapper/events"
	"github.com/dapperlabs/dapper-flow-api/pkg/retry"
	timeutils "github.com/dapperlabs/dapper-flow-api/pkg/time"
)

// BackoffMiddleware applies a retry backoff to an event handler.
type BackoffMiddleware struct {
	tracker *retry.BackoffTracker
}

// NewBackoffMiddleware returns a new BackoffMiddleware.
func NewBackoffMiddleware(b *retry.BackoffTracker) *BackoffMiddleware {
	return &BackoffMiddleware{b}
}

// HandleEvent calls the `next` event handler and backs off if the handler retries.
func (m *BackoffMiddleware) HandleEvent(
	ctx context.Context,
	eventData []byte,
	next EventHandlerFunc,
) *eventclient.EventProcessResult {
	event := &events.AnyEvent{}
	if err := proto.Unmarshal(eventData, event); err == nil {
		eventID := event.GetHeader().GetEventId()
		result := next(ctx, eventData)
		if eventID != "" {
			if result.Retry {
				// WaitUntilDoneWithTimeout until the backoff delay is over, or until the context cancels
				timeutils.WaitUntilDoneWithTimeout(ctx, m.tracker.Inc(eventID))
			} else {
				m.tracker.Reset(eventID)
			}
		}

		return result
	}

	return next(ctx, eventData)
}
