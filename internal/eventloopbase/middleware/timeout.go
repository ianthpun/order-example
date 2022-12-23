package middleware

import (
	"context"
	"time"

	"github.com/dapperlabs/dibs/v2/eventclient"
)

// TimeoutMiddleware applies a timeout to the event handler.
type TimeoutMiddleware struct {
	Timeout time.Duration
}

// NewTimeoutMiddleware returns a new TimeoutMiddleware.
func NewTimeoutMiddleware(timeout time.Duration) *TimeoutMiddleware {
	return &TimeoutMiddleware{timeout}
}

// HandleEvent calls the `next` event handler with a context timeout.
func (m *TimeoutMiddleware) HandleEvent(
	ctx context.Context,
	eventData []byte,
	next EventHandlerFunc,
) *eventclient.EventProcessResult {
	ctx, cancel := context.WithTimeout(ctx, m.Timeout)
	defer cancel()
	return next(ctx, eventData)
}
