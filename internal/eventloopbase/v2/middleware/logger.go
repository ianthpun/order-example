package middleware

import (
	"context"

	"github.com/dapperlabs/dibs/v3/eventclient"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"

	"github.com/dapperlabs/dapper-flow-api/internal/services"
	"github.com/dapperlabs/dapper-flow-api/pkg/protobuf/dapper/events"
)

// LoggerMiddleware applies a logger to the event handler. The logger will log the result returned by the event handler.
type LoggerMiddleware struct {
	logger *services.Logger
}

// NewLoggerMiddleware returns a new LoggerMiddleware.
func NewLoggerMiddleware(logger *services.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{logger}
}

// HandleEvent calls the `next` event handler logs the result.
func (m *LoggerMiddleware) HandleEvent(
	ctx context.Context,
	eventData []byte,
	next EventHandlerFunc,
) *eventclient.EventProcessResult {
	eventName := "Unknown"

	// Try get the event name
	event := &events.AnyEvent{}
	if err := proto.Unmarshal(eventData, event); err != nil {
		m.logger.With(zap.Error(err)).
			Warn("unable to unmarshal generic event")
	} else {
		// This event has a header, so we can use it to get the name of the event
		eventName = event.GetHeader().GetEventName()
	}

	logger := m.logger.With(
		zap.String("event_id", event.GetHeader().GetEventId()),
		zap.String("request_id", event.GetHeader().GetRequestId()),
		zap.String("event_name", eventName),
	)
	logger.Debug("handling event")

	// add logger to context
	ctx = ctxzap.ToContext(ctx, logger.Logger)

	result := next(ctx, eventData)

	logger = logger.With(
		zap.Bool("retry", result.Retry),
		zap.Bool("failed", result.Failed),
	)

	if result.Failed || result.Retry {
		logger.Error(result.Message)
	} else {
		logger.Debug(result.Message)
	}

	return result
}
