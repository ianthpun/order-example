package eventclient

import "context"

type Subscriber interface {
	Subscribe(ctx context.Context, topic string) (<-chan *Message, error)
	Close() error
}

type Publisher interface {
	Publish(topic string, messages ...*Message) error
	Close() error
}
