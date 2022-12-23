package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapperlabs/dapper-flow-api/pkg/protobuf/dapper/events"
	"github.com/dapperlabs/dibs/v2/eventclient"
	"github.com/golang/protobuf/proto"
)

var mid *PrometheusMiddleware

func middleware() (*PrometheusMiddleware, error) {
	if mid != nil {
		return mid, nil
	}

	var err error
	mid, err = NewPrometheusMiddleware()
	return mid, err
}

func TestPrometheusMiddleware_HandleEvent(t *testing.T) {
	mid, err := middleware()
	require.Nil(t, err)

	event := &events.AnyEvent{
		Header: &events.Header{
			EventName: "EventName",
			EventId:   "EventId",
			RequestId: "RequestId",
		},
	}
	eventBytes, err := proto.Marshal(event)
	require.Nil(t, err)

	handler := func(ctx context.Context, eventData []byte) *eventclient.EventProcessResult {
		event := &events.AnyEvent{}
		err := proto.Unmarshal(eventData, event)
		require.Nil(t, err)

		header := event.GetHeader()

		return &eventclient.EventProcessResult{
			Retry:   false,
			Failed:  false,
			Message: fmt.Sprintf("%s %s %s", header.GetEventId(), header.GetEventName(), header.GetRequestId()),
		}
	}

	res := mid.HandleEvent(context.Background(), eventBytes, handler)
	require.Equal(t, "EventId EventName RequestId", res.Message)

	// Make sure the right metrics are emitted
	resp := doRequest(promhttp.Handler(), http.MethodGet, "/", "")
	require.Equal(t, http.StatusOK, resp.Code)

	body := resp.Body.String()
	assert.Contains(t, body, `events_handled_total{failed="false",name="EventName",retry="false"} 1`)
	assert.Contains(t, body, `event_handler_duration_ms_bucket_count{name="EventName"} 1`)
}

func TestPrometheusMiddleware_HandleEventWithUnknownEvent(t *testing.T) {
	mid, err := middleware()
	require.Nil(t, err)

	handler := func(ctx context.Context, eventData []byte) *eventclient.EventProcessResult {
		return &eventclient.EventProcessResult{
			Message: "test",
		}
	}

	res := mid.HandleEvent(context.Background(), make([]byte, 5), handler)
	require.Equal(t, "test", res.Message)

	// Make sure the right metrics are emitted
	resp := doRequest(promhttp.Handler(), http.MethodGet, "/", "")
	require.Equal(t, http.StatusOK, resp.Code)

	body := resp.Body.String()
	assert.Contains(t, body, `events_handled_total{failed="false",name="Unknown",retry="false"} 1`)
	assert.Contains(t, body, `event_handler_duration_ms_bucket_count{name="Unknown"} 1`)
}

func doRequest(handler http.Handler, method string, target string, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, strings.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	return w
}
