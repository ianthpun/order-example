package watermill

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	"order-sample/pkg/eventclient"
)

var _ eventclient.Router = (*router)(nil)

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
	subscriber eventclient.Subscriber,
	publishTopic string,
	publisher eventclient.Publisher,
	handlerFunc eventclient.HandlerFunc,
	middlewares ...eventclient.MiddlewareFunc,
) error {
	handler := r.watermillRouter.AddHandler(
		handlerName,
		subscribeTopic,
		NewSubscriberAdapter(subscriber),
		publishTopic,
		NewPublisherAdapter(publisher),
		ToWatermillHandlerFunc(handlerFunc),
	)

	for _, md := range middlewares {
		handler.AddMiddleware(ToWatermillMiddlewareFunc(md))
	}

	return nil
}

func (*router) RegisterWithoutPublishing(
	handlerName string,
	subscribeTopic string,
	subscriber eventclient.Subscriber,
	publishTopic string,
	publisher eventclient.Publisher,
	handlerFunc eventclient.HandlerFunc,
	middlewares ...eventclient.MiddlewareFunc,
) error {
	//TODO implement me
	panic("implement me")
}

func (r *router) AddMiddleware(middlewares ...eventclient.MiddlewareFunc) {
	// TODO: add this here
}

func (r *router) Run(ctx context.Context) error {
	return r.watermillRouter.Run(ctx)
}
