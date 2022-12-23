package middleware

import (
	"context"
	"fmt"

	"github.com/dapperlabs/dibs/v3/eventclient"
	"github.com/golang/protobuf/proto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/dapperlabs/dapper-flow-api/pkg/protobuf/dapper/events"
)

// NopTracingMiddleware is a tracing event loop middleware that calls the next func,
// and does nothing else.
type NopTracingMiddleware struct{}

// HandleEvent implements Middleware and calls next
func (m *NopTracingMiddleware) HandleEvent(
	ctx context.Context,
	eventData []byte,
	next EventHandlerFunc,
) *eventclient.EventProcessResult {
	return next(ctx, eventData)
}

// EventTracingMiddleware starts a new span for the event consumer.
type EventTracingMiddleware struct {
	topic               string
	consumerGroupID     string
	consumerContentType string
}

// NewEventTracingMiddleware returns a new EventTracingMiddleware.
func NewEventTracingMiddleware(topic string, conf eventclient.Config) *EventTracingMiddleware {
	return &EventTracingMiddleware{
		topic:               topic,
		consumerGroupID:     conf.GroupID,
		consumerContentType: string(conf.KafkaConsumerConfig.ContentType),
	}
}

// HandleEvent calls the `next` event handler and backs off if the handler retries.
func (m *EventTracingMiddleware) HandleEvent(
	ctx context.Context,
	eventData []byte,
	next EventHandlerFunc,
) *eventclient.EventProcessResult {
	eventName := "Unknown"
	event := &events.AnyEvent{}

	if err := proto.Unmarshal(eventData, event); err == nil {
		eventName = event.GetHeader().GetEventName()
	}

	tracer := otel.Tracer("eventloop-kafka-consumer")
	attrs := []attribute.KeyValue{
		attribute.String("messaging.system", "kafka"),
		attribute.String("messaging.operation", "receive"),
		attribute.String("messaging.destination", m.topic),
		attribute.String("messaging.destination_kind", "topic"),
		attribute.String("messaging.kafka.message_key", ""),
		attribute.String("messaging.kafka.consumer.group_id", m.consumerGroupID),
		attribute.String("messaging.kafka.client_id", ""),
		attribute.String("messaging.kafka.partition", ""),
		attribute.String("messaging.kafka.consumer.content_type", m.consumerContentType),
		attribute.String("messaging.event.event_id", event.GetHeader().GetEventId()),
		attribute.String("messaging.event.name", eventName),
		attribute.String("messaging.event.request_id", event.GetHeader().GetRequestId()),
	}

	var opts []trace.SpanStartOption
	opts = append(opts,
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(attrs...),
	)

	ctx, span := tracer.Start(ctx, fmt.Sprintf("%s handler", m.topic), opts...)
	defer span.End()

	result := next(ctx, eventData)
	span.SetAttributes(
		attribute.Bool("messaging.event.retry", result.Retry),
		attribute.Bool("messaging.event.failed", result.Failed),
		attribute.String("messaging.event.message", result.Message),
	)

	for _, outboxMessage := range result.Outbox {
		subEventName := "unknown"
		// Here we are creating a new sub span for each outbox event that was publish from the
		// handling of the parent span so that we can search by the outbox event name.
		// We don't have full access to the eventclient publishing from here so these sub spans won't be propagated
		// and be linked up on the event consumer side so they may end up not being that useful, in which we can turn these
		// into span events.
		_, subSpan := tracer.Start(ctx, outboxMessage.Topic, trace.WithSpanKind(trace.SpanKindProducer))

		if messageType := proto.MessageReflect(outboxMessage.Message.Proto); messageType != nil {
			subEventName = messageType.Type().Descriptor().Syntax().String()
		}

		subSpan.SetAttributes(
			attribute.String("messaging.event.name", subEventName),
			attribute.String("messaging.event.reason", outboxMessage.Reason),
			attribute.String("messaging.system", "kafka"),
			attribute.String("messaging.operation", "process"),
			attribute.String("messaging.destination", outboxMessage.Topic),
			attribute.String("messaging.destination_kind", "topic"),
		)
		subSpan.End()
	}

	return result
}
