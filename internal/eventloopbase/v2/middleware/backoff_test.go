package middleware

import (
	"context"
	"testing"
	"time"

	"github.com/dapperlabs/dibs/v3/eventclient"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"github.com/dapperlabs/dapper-flow-api/pkg/protobuf/dapper/events"
	"github.com/dapperlabs/dapper-flow-api/pkg/retry"
)

func TestBackoffMiddleware(t *testing.T) {
	mid := NewBackoffMiddleware(retry.NewBackoffTracker(time.Millisecond*10, time.Millisecond*100, 2))

	event := &events.AnyEvent{
		Header: &events.Header{
			EventName: "EventName",
			EventId:   "EventId",
			RequestId: "RequestId",
		},
	}
	eventBytes, err := proto.Marshal(event)
	assert.Nil(t, err)

	backoffTimes := []time.Duration{
		time.Millisecond * 10,
		time.Millisecond * 20,
		time.Millisecond * 40,
		time.Millisecond * 80,
		time.Millisecond * 100, // Max backoff reached here
		time.Millisecond * 100,
		time.Millisecond * 100,
	}

	retryCount := 0
	handler := func(ctx context.Context, eventData []byte) *eventclient.EventProcessResult {
		// Stop retrying after all backoff times have been hit
		shouldRetry := retryCount < len(backoffTimes)
		retryCount++
		return &eventclient.EventProcessResult{
			Retry: shouldRetry,
		}
	}

	// Check that the middleware applies backoff sleeps correctly
	for _, minExpected := range backoffTimes {
		start := time.Now()
		mid.HandleEvent(context.Background(), eventBytes, handler)
		assert.True(t, time.Now().Sub(start) >= minExpected)
		println(time.Now().Sub(start).String())
	}

	// Check that the middleware resets the backoff correctly
	start := time.Now()
	mid.HandleEvent(context.Background(), eventBytes, handler)
	assert.True(t, time.Now().Sub(start) <= time.Millisecond*10)
}
