package eventloopbase

import (
	"context"
	"testing"

	"github.com/dapperlabs/dibs/v2/eventclient"
	"github.com/go-pg/pg/v9"

	"github.com/dapperlabs/dapper-flow-api/internal/eventloopbase/middleware"
)

type testMid struct {
	value string
}

func (t *testMid) HandleEvent(
	ctx context.Context,
	eventData []byte,
	next middleware.EventHandlerFunc,
) *eventclient.EventProcessResult {
	return next(ctx, append(eventData, []byte(t.value)...))
}

func TestMiddleware(t *testing.T) {
	b := New(nil, nil, eventclient.Config{}, &testMid{"1"}, &testMid{"2"})
	b.Register("test", func(ctx context.Context, eventData []byte) *eventclient.EventProcessResult {
		return Success(string(eventData))
	}, &testMid{"3"})

	// There should only be 2 default middleware
	if len(b.middleware) != 2 {
		t.FailNow()
	}

	// All middleware should be applied
	res := b.eventHandlers["test"](context.Background(), &pg.Tx{}, []byte("something"))
	expected := "something321"
	if res.Message != expected {
		t.Logf("expected %s, but got %s", expected, res.Message)
		t.FailNow()
	}
}
