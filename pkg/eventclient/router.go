package eventclient

import (
	"context"
)

type Router interface {
	Register(
		handlerName string,
		subscribeTopic string,
		subscriber Subscriber,
		publishTopic string,
		publisher Publisher,
		handlerFunc HandlerFunc,
		middlewares ...MiddlewareFunc,
	) error
	RegisterWithoutPublishing(
		handlerName string,
		subscribeTopic string,
		subscriber Subscriber,
		publishTopic string,
		publisher Publisher,
		handlerFunc HandlerFunc,
		middlewares ...MiddlewareFunc,
	) error
	AddMiddleware(
		middlewares ...MiddlewareFunc,
	)
	Run(ctx context.Context) error
}

type HandlerFunc func(msg *Message) ([]*Message, error)
type MiddlewareFunc func(handlerFunc HandlerFunc) HandlerFunc
