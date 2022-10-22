package kafka

import (
	"crypto/tls"
	"github.com/Shopify/sarama"
)

func WithTLS(tls *tls.Config) func(*sarama.Config) {
	return func(conf *sarama.Config) {
		conf.Net.TLS.Enable = true
		conf.Net.TLS.Config = tls
	}
}

func WithSASL(username, password string) func(*sarama.Config) {
	return func(conf *sarama.Config) {
		conf.Net.SASL.Enable = true
		conf.Net.SASL.User = username
		conf.Net.SASL.Password = password
	}
}
