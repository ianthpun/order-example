package kafka

import (
	"crypto/tls"
	"github.com/Shopify/sarama"
	pkgkafka "github.com/ThreeDotsLabs/watermill-kafka/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	eventclient2 "order-sample/pkg/eventclient"
	"order-sample/pkg/eventclient/watermill"
)

type Publisher struct {
	pub message.Publisher
}

var _ eventclient2.Publisher = (*Publisher)(nil)

type PublisherConfig struct {
	Username string
	Password string
	TLS      *tls.Config
}

func NewPublisher(
	brokers []string,
	options ...func(config *sarama.Config),
) (*Publisher, error) {
	saramaConf := pkgkafka.DefaultSaramaSyncPublisherConfig()

	for _, o := range options {
		o(saramaConf)
	}

	pub, err := pkgkafka.NewPublisher(
		brokers,
		pkgkafka.NewWithPartitioningMarshaler(func(topic string, msg *message.Message) (string, error) {
			return msg.Metadata.Get(eventclient2.PartitionKey), nil
		}),
		saramaConf,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		pub: pub,
	}, nil
}

func (p *Publisher) Publish(topic string, messages ...*eventclient2.Message) error {
	return p.pub.Publish(topic, watermill.MessagesToWatermillMessages(messages)...)
}

func (p *Publisher) Close() error {
	return p.pub.Close()
}
