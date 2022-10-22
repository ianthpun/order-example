package kafka

import (
	"context"
	"crypto/tls"
	"github.com/Shopify/sarama"
	pkgkafka "github.com/ThreeDotsLabs/watermill-kafka/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	eventclient2 "order-sample/pkg/eventclient"
)

type Subscriber struct {
	sub message.Subscriber
}

var _ eventclient2.Subscriber = (*Subscriber)(nil)

type SubscriberConfig struct {
	Username string
	Password string
	TLS      *tls.Config
}

func NewSubscriber(
	brokers []string,
	consumerGroup string,
	options ...func(config *sarama.Config),
) (*Subscriber, error) {
	saramaConf := pkgkafka.DefaultSaramaSubscriberConfig()

	for _, o := range options {
		o(saramaConf)
	}

	sub, err := pkgkafka.NewSubscriber(
		pkgkafka.SubscriberConfig{
			Brokers:       brokers,
			ConsumerGroup: consumerGroup,
		},
		saramaConf,
		pkgkafka.NewWithPartitioningMarshaler(func(topic string, msg *message.Message) (string, error) {
			return msg.Metadata.Get(eventclient2.PartitionKey), nil
		}),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		sub: sub,
	}, nil
}
func (s Subscriber) Subscribe(ctx context.Context, topic string) (<-chan *eventclient2.Message, error) {
	//TODO implement me
	panic("implement me")
}

func (s Subscriber) Close() error {
	//TODO implement me
	panic("implement me")
}
