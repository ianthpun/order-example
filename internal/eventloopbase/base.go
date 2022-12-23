package eventloopbase

import (
	"context"
	"errors"
	"fmt"
	"order-sample/internal/eventloopbase/middleware"
	"strings"
	"time"

	"github.com/dapperlabs/dibs/v2/eventclient"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Middleware is any function that wraps an event handler.
type Middleware interface {
	HandleEvent(ctx context.Context, eventData []byte, next middleware.EventHandlerFunc) *eventclient.EventProcessResult
}

// EventHeader is an interface that represents an event header.
type EventHeader interface {
	GetEventId() string
	GetEventName() string
	GetRequestId() string
}

// BaseEventLoop implements basic eventloop functionality.
type BaseEventLoop struct {
	KafkaClient eventclient.Provider
	eventclient.Config

	logger            *zap.Logger
	eventHandlers     map[string]func(context.Context, interface{}, []byte) *eventclient.EventProcessResult
	errHandlers       map[string]func(error)
	defaultErrHandler func(error)
	middleware        []Middleware
}

// New creates a new BaseEventLoop.
func New(
	logger *zap.Logger,
	kafkaClient eventclient.Provider,
	ecConfig eventclient.Config,
	middleware ...Middleware,
) *BaseEventLoop {
	b := &BaseEventLoop{
		KafkaClient:   kafkaClient,
		Config:        ecConfig,
		logger:        logger,
		eventHandlers: make(map[string]func(context.Context, interface{}, []byte) *eventclient.EventProcessResult),
		errHandlers:   make(map[string]func(error)),
		middleware:    middleware,
	}
	b.SetDefaultErrHandler(b.handleErrorFunc)
	return b
}

// SetDefaultErrHandler sets the default error handler that is called by eventclient when an event handler.
// returns an error.
func (b *BaseEventLoop) SetDefaultErrHandler(handler func(error)) {
	b.defaultErrHandler = handler
}

// Register registers an event handler to be called when an event is received from the given topic.
// The default error handler will be called if the handler returns an error. All default middleware registered on this
// BaseEventLoop will be applied to the given handler. Any additional middleware passed to this function will be
// applied after the default middleware is applied.
func (b *BaseEventLoop) Register(topic string, eventHandler middleware.EventHandlerFunc, mid ...Middleware) {
	eventLoopMid := append(b.middleware, mid...)

	if tracingMid, ok := RegisterEventLoopMid(topic, b.Config).(Middleware); ok {
		// tracing middleware should be first, so we always create a trace
		eventLoopMid = append([]Middleware{tracingMid}, eventLoopMid...)
	}

	// Wrap with event handler with default middleware, and then with the provided middleware
	wrappedHandler := applyMiddleware(eventLoopMid, eventHandler)

	b.RegisterWithErrHandler(topic, wrappedHandler, b.defaultErrHandler)
}

// RegisterWithErrHandler registers an event handler to be called when an event is received from the given topic.
// The given error handler will be called if the handler returns an error.
func (b *BaseEventLoop) RegisterWithErrHandler(
	topic string,
	eventHandler middleware.EventHandlerFunc,
	errHandler func(error),
) {
	// Register the event handler, wrapping it in a function that type casts the tx provided by eventclient
	b.eventHandlers[topic] = wrap(eventHandler)

	// Register an error handler
	b.errHandlers[topic] = errHandler
}

// Start starts an event loop for each registered topic/
func (b *BaseEventLoop) Start(ctx context.Context) error {
	for topic, handler := range b.eventHandlers {
		b.logger.Info("Starting eventloop for topic", zap.String("topic", topic))
		if err := b.KafkaClient.StartEventLoop(ctx, topic, handler, b.errHandlers[topic]); err != nil {
			b.logger.Error("Failed to start event loop for topic", zap.String("topic", topic), zap.Error(err))

			return err
		}

		b.logger.Info("Successfully registering event loop for topic", zap.String("topic", topic))
	}
	return nil
}

// Fail is a utility function that returns an EventProcessResult containing the formatted failure message with no outbox.
func Fail(format string, args ...interface{}) *eventclient.EventProcessResult {
	return &eventclient.EventProcessResult{
		Failed:  true,
		Message: fmt.Sprintf(format, args...),
		Outbox:  nil,
	}
}

// Success is a utility function that returns an EventProcessResult containing the formatted success message with no outbox.
func Success(format string, args ...interface{}) *eventclient.EventProcessResult {
	return &eventclient.EventProcessResult{
		Message: fmt.Sprintf(format, args...),
	}
}

// Empty is a utility function that returns an EventProcessResult with no outbox
func Empty() *eventclient.EventProcessResult {
	return &eventclient.EventProcessResult{
		Failed: false,
		Outbox: nil,
	}
}

// Retry is a utility function that returns an EventProcessResult with Retry flag set to true. This will
// cause the event to be retried indefinitely
func Retry() *eventclient.EventProcessResult {
	return &eventclient.EventProcessResult{
		Retry: true,
	}
}

// HandleErr is a utility function that handles an error return from an external call by either retrying or failing
// the event and moving on.
func HandleErr(
	err error,
	msgFormat string,
	msgArgs ...interface{},
) *eventclient.EventProcessResult {
	if IsTransient(err) {
		return Retry()
	}

	errMsg := fmt.Sprintf(" : %v", err)

	return Fail(msgFormat+errMsg, msgArgs...)
}

// RetryTransientErrors retries function f as long as a transient error is returned for maxDuration.
// minWaitDuration is the minimum duration between each retry. It errors if it isn't greater than 0.
func RetryTransientErrors(
	ctx context.Context,
	maxDuration time.Duration,
	minWaitDuration time.Duration,
	f func() error,
) (err error) {
	if minWaitDuration <= 0 {
		return errors.New("minWaitDuration must be greater than 0")
	}

	maxDurationTimer := time.NewTimer(maxDuration)
	defer maxDurationTimer.Stop()

	ticker := time.NewTicker(minWaitDuration)
	defer ticker.Stop()

	for {
		if err = f(); err == nil || !IsTransient(err) {
			return
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-maxDurationTimer.C:
			return
		case <-ticker.C:
			continue
		}
	}
}

// IsTransient checks if the error is transient. Transient errors are often retried by event handlers
// until they disappear (hopefully).
func IsTransient(err error) bool {
	errCode := status.Code(err)
	if errCode == codes.Unavailable || errCode == codes.DeadlineExceeded {
		// Temporary gRPC error.
		return true
	}

	// check for context deadline exceeded if gRPC error is wrapped. flow-go-sdk wraps the error before returning.
	// this is hacked as string check for now, since we can't check grpc status on wrapped errors
	// more info https://github.com/grpc/grpc-go/issues/2934
	// flow issue https://github.com/onflow/flow-go-sdk/issues/32
	if strings.Contains(err.Error(), "context deadline exceeded") {
		// Temporary gRPC error.
		return true
	}

	// check for connection pool timeouts
	if strings.HasSuffix(err.Error(), "connection pool timeout") {
		// a database connection attempt was timed out because the pool is too busy
		// this is considered transient as we assume pool connections eventually become available
		return true
	}

	// check for transient dal EOF errors
	if strings.Contains(err.Error(), "dal: network error: EOF") {
		return true
	}

	// check for transient EOF error on Waterhose
	if strings.Contains(err.Error(), "Internal") &&
		(strings.Contains(err.Error(), "InsertTransaction: EOF")) {
		return true
	}

	// transient KMS key generation error
	if strings.Contains(err.Error(), "Internal") &&
		(strings.Contains(err.Error(), "GetPublicKeyWithKeyConfig: %!s(<nil>)")) {
		return true
	}

	// check for KMS rate limit errors
	if strings.Contains(err.Error(), "ResourceExhausted") &&
		(strings.Contains(err.Error(), "Quota exceeded for quota metric") ||
			strings.Contains(err.Error(), "exceeded limit for metric")) {
		return true
	}

	return false
}

// Logger returns the base logger.
func (b *BaseEventLoop) Logger() *zap.Logger {
	return b.logger
}

// LoggerFor returns a structured field logger for an event. The logger will include the event name and ID as fields
// on each log message.
func (b *BaseEventLoop) LoggerFor(header EventHeader) *zap.Logger {
	return b.logger.With(
		zap.String("event_name", header.GetEventName()),
		zap.String("event_id", header.GetEventId()),
		zap.String("request_id", header.GetRequestId()),
	)
}

// handleErrorFunc used for handling errors coming from the event loop
func (b *BaseEventLoop) handleErrorFunc(err error) {
	// log ErrProcessFuncRetry error as Debug log without stack trace
	if errors.Is(err, eventclient.ErrProcessFuncRetry) {
		b.logger.Debug("event handler requested retry")
		return
	}

	b.logger.Error("error in event loop handler", zap.Error(err))
}

// wrap wraps the given handler in a handler that is compatible with eventclient.
// This wrapper ensures that we're passing a *pg.Tx to the handler if persistence is on,
// and that the transaction is also available in the context that is passed to the handler so our DAL methods
// can use it if need be.
func wrap(handler middleware.EventHandlerFunc) func(
	context.Context,
	interface{},
	[]byte,
) *eventclient.EventProcessResult {
	return func(ctx context.Context, tx interface{}, bytes []byte) *eventclient.EventProcessResult {
		return handler(ctx, bytes)
	}
}

// applyMiddleware wraps the given handler in the middleware, returning the new wrapped handler.
func applyMiddleware(middleware []Middleware, handler middleware.EventHandlerFunc) middleware.EventHandlerFunc {
	if len(middleware) == 0 {
		return handler
	}

	// Wrap handler with next middleware
	wrappedHandler := func(ctx context.Context, bytes []byte) *eventclient.EventProcessResult {
		return middleware[0].HandleEvent(ctx, bytes, handler)
	}

	// Wrap handler with remaining middleware
	return applyMiddleware(middleware[1:], wrappedHandler)
}

// RegisterEventLoopMid will register an eventloop middleware
func RegisterEventLoopMid(topic string, config eventclient.Config) interface{} {
	return middleware.NopTracingMiddleware{}
}
