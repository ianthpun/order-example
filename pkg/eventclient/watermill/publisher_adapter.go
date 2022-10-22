package watermill

import (
	pkgwatermill "github.com/ThreeDotsLabs/watermill/message"
	"order-sample/pkg/eventclient"
)

var _ pkgwatermill.Publisher = (*publisherAdapter)(nil)

type publisherAdapter struct {
	publisher eventclient.Publisher
}

func NewPublisherAdapter(p eventclient.Publisher) *publisherAdapter {
	return &publisherAdapter{}
}

func (w publisherAdapter) Publish(topic string, messages ...*pkgwatermill.Message) error {
	return w.publisher.Publish(topic, WatermillMessagesToMessages(messages)...)
}

func (w publisherAdapter) Close() error {
	return w.publisher.Close()
}
