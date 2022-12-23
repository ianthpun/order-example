package middleware

import (
	"context"

	"github.com/dapperlabs/dibs/v2/eventclient"
)

// EventHandlerFunc is any function that handles events.
type EventHandlerFunc func(context.Context, []byte) *eventclient.EventProcessResult
