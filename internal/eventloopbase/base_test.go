package eventloopbase

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	dalerrs "github.com/dapperlabs/dapper-flow-api/internal/dal/errors"
	"github.com/dapperlabs/dapper-flow-api/internal/services"
)

func TestIsTransient(t *testing.T) {
	t.Run("timeouts", func(t *testing.T) {
		// create context deadline error
		ctx, cancelFunc := context.WithDeadline(context.Background(), time.Now())
		defer cancelFunc()
		_, err := grpc.DialContext(
			ctx,
			"localhost:8207",
			services.DefaultDialOpts...,
		)

		// check unwrapped
		assert.True(t, IsTransient(err))

		// check wrapped
		err = fmt.Errorf("client: %w", err)
		assert.True(t, IsTransient(err))
	})

	t.Run("Postgres errors", func(t *testing.T) {
		assert.True(t, IsTransient(dalerrs.ErrTransient))
	})
}

func TestRetryTransientErrors(t *testing.T) {
	var (
		ctx             = context.Background()
		testDuration    = 3 * time.Second
		minWaitDuration = 100 * time.Millisecond
	)

	t.Run("on transient errors it retries until testDuration is reached", func(t *testing.T) {
		start := time.Now()
		errTransient := errors.New("context deadline exceeded")
		assert.True(t, IsTransient(errTransient))

		retryErr := RetryTransientErrors(ctx, testDuration, minWaitDuration, func() error {
			return errTransient
		})

		assert.True(t, IsTransient(retryErr))
		assert.GreaterOrEqual(t, time.Now().Sub(start), testDuration)
	})

	t.Run("on non-transient errors it returns the error after the first try", func(t *testing.T) {
		start := time.Now()
		errNonTransient := errors.New("not transient")
		callsCount := 0
		retryErr := RetryTransientErrors(ctx, testDuration, minWaitDuration, func() error {
			callsCount += 1
			return errNonTransient
		})
		assert.False(t, IsTransient(retryErr))
		// should terminate instantly
		assert.Equal(t, 1, callsCount)
		assert.LessOrEqual(t, time.Now().Sub(start), testDuration)
	})

	t.Run("when minWaitDuration is set it pauses between retries", func(t *testing.T) {
		start := time.Now()
		errTransient := errors.New("context deadline exceeded")
		assert.True(t, IsTransient(errTransient))
		callsCount := 0

		retryErr := RetryTransientErrors(ctx, testDuration, 500*time.Millisecond, func() error {
			callsCount += 1
			return errTransient
		})

		assert.True(t, IsTransient(retryErr))
		assert.Equal(t, 6, callsCount)
		assert.GreaterOrEqual(t, time.Now().Sub(start), testDuration)
	})

	t.Run("when context timeout is reached", func(t *testing.T) {
		errTransient := errors.New("context deadline exceeded")
		assert.True(t, IsTransient(errTransient))

		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		callsCount := 0
		defer cancel()

		retryErr := RetryTransientErrors(ctx, testDuration, 600*time.Millisecond, func() error {
			callsCount += 1
			return errTransient
		})

		assert.Equal(t, 2, callsCount)
		assert.Error(t, retryErr)
		assert.Equal(t, context.DeadlineExceeded, retryErr)
	})
}
