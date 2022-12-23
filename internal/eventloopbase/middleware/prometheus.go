package middleware

import (
	"context"
	"strconv"
	"time"

	"github.com/dapperlabs/dibs/v2/eventclient"
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/dapperlabs/dapper-flow-api/pkg/protobuf/dapper/events"
)

var (
	eventHandledCounterLabels = []string{"name", "failed", "retry"}
	timeToHandleHistLabels    = []string{"name"}
)

// PrometheusMiddleware emits event metrics.
type PrometheusMiddleware struct {
	eventHandledCounter *prometheus.CounterVec
	timeToHandleHist    *prometheus.HistogramVec
}

// NewPrometheusMiddleware returns a new PrometheusMiddleware.
// WARNING: This function will return an error if it is called multiple times in the same process because you
// cannot register the same metrics more than once.
func NewPrometheusMiddleware() (*PrometheusMiddleware, error) {
	// Create and register a metric for counting the number of events handled with labels indicating the event name,
	// whether the handler will retry, and whether the handler failed.
	eventHandledCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "events_handled_total",
			Help: "Total number of events handled.",
		},
		eventHandledCounterLabels,
	)
	if err := prometheus.Register(eventHandledCounter); err != nil {
		return nil, err
	}

	// Create and register a metric for tracking the amount of time it takes to handle the event with a label indicating
	// the event name. Samples will be placed into buckets of the following sizes (in milliseconds):
	// 50, 100, 200, 400, 800, 1600, 3200, 6400, 12800, 25600.
	timeToHandleHist := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "event_handler_duration_ms_bucket",
			Help:    "The time taken to handle an event.",
			Buckets: prometheus.ExponentialBuckets(50, 2, 10),
		},
		timeToHandleHistLabels,
	)
	if err := prometheus.Register(timeToHandleHist); err != nil {
		return nil, err
	}

	return &PrometheusMiddleware{
		eventHandledCounter: eventHandledCounter,
		timeToHandleHist:    timeToHandleHist,
	}, nil
}

// HandleEvent calls the `next` handler and emits metrics about the event.
func (m *PrometheusMiddleware) HandleEvent(
	ctx context.Context,
	eventData []byte,
	next EventHandlerFunc,
) *eventclient.EventProcessResult {
	// Try get the event name
	eventName := "Unknown"
	event := &events.AnyEvent{}
	if err := proto.Unmarshal(eventData, event); err == nil {
		// This event has a header, so we can use it to get the name of the event
		eventName = event.GetHeader().GetEventName()
	}

	// Time the event handler
	start := time.Now()

	// Call event handler
	res := next(ctx, eventData)

	// Update metrics
	m.eventHandledCounter.WithLabelValues(
		eventName,
		strconv.FormatBool(res.Failed),
		strconv.FormatBool(res.Retry),
	).Inc()
	m.timeToHandleHist.WithLabelValues(eventName).Observe(float64(time.Since(start).Milliseconds()))

	return res
}
