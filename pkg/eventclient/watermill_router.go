package eventclient

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	watermill2 "order-sample/pkg/eventclient/watermill"
)

var _ Router = (*router)(nil)

type router struct {
	watermillRouter *message.Router
}

func NewRouter() *router {
	r, err := message.NewRouter(
		message.RouterConfig{},
		nil,
	)
	if err != nil {
		panic(err)
	}

	return &router{
		watermillRouter: r,
	}
}

func (r *router) Register(
	handlerName string,
	subscribeTopic string,
	subscriber Subscriber,
	publishTopic string,
	publisher Publisher,
	handlerFunc HandlerFunc,
	middlewares ...MiddlewareFunc,
) error {
	handler := r.watermillRouter.AddHandler(
		handlerName,
		subscribeTopic,
		watermill2.NewSubscriberAdapter(subscriber),
		publishTopic,
		watermill2.NewPublisherAdapter(publisher),
		watermill2.ToWatermillHandlerFunc(handlerFunc),
	)

	for _, md := range middlewares {
		handler.AddMiddleware(watermill2.ToWatermillMiddlewareFunc(md))
	}

	return nil
}

func (*router) RegisterWithoutPublishing(
	handlerName string,
	subscribeTopic string,
	subscriber Subscriber,
	publishTopic string,
	publisher Publisher,
	handlerFunc HandlerFunc,
	middlewares ...MiddlewareFunc,
) error {
	//TODO implement me
	panic("implement me")
}

func (r *router) AddMiddleware(middlewares ...MiddlewareFunc) {
	// TODO: add this here
}

func (r *router) Run(ctx context.Context) error {
	return r.watermillRouter.Run(ctx)
}
