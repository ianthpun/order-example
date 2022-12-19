package watermill

import (
	"context"
	pkgwatermill "github.com/ThreeDotsLabs/watermill/message"
	"order-sample/pkg/eventclient"
)

type subscriberAdapter struct {
	subscriber eventclient.Subscriber
}

func NewSubscriberAdapter(p eventclient.Subscriber) *subscriberAdapter {
	return &subscriberAdapter{}
}

func (w subscriberAdapter) Subscribe(ctx context.Context, topic string) (<-chan *pkgwatermill.Message, error) {
	//return w.subscriber.Subscribe(ctx, topic)
	return nil, nil
}

func (w subscriberAdapter) Close() error {
	return w.Close()
}
