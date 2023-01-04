package config

import (
	"crypto/tls"
	kafka "github.com/dapperlabs/dibs/v2/multikafka"
	"order-sample/internal/services"
	"order-sample/pkg/temporal"
	"time"
)

type (
	// BaseConfig for the standard configuration of the service
	BaseConfig = services.BaseConfig
	// GRPCServiceConfig for GRPC services
	GRPCServiceConfig = services.GRPCServiceConfig
	// KafkaConfig for provider configuration
	KafkaConfig = kafka.Config
	// TemporalConfig is the config for Temporal
	TemporalConfig = temporal.Config
	// EventloopConfig for eventloop
	EventloopConfig = services.EventLoopConfig
)

type Config struct {
	// Base service config
	BaseConfig

	// Kafka config
	KafkaConfig
	KafkaConsumerMaxWait time.Duration `default:"500ms"`

	// Eventloop config
	EventloopConfig

	// GRPC Config
	GRPCServiceConfig

	// Temporal Config
	TemporalConfig

	RunMigrations bool `default:"false"`
}

// ComputeDependents computes variables that depend on previous config vars
// and it performs sanity checks on convenient defaults
func (c *Config) ComputeDependents() (err error) {
	// Select a Kafka broker. We do this because in some cases there are multiple space-separated brokers
	// in the KafkaBroker config variable.
	c.KafkaPrimaryBroker = c.BrokerPrimary()
	c.KafkaSecondaryBroker = c.BrokerSecondary()

	// Compute TLS configs if no migrations are being run
	if !c.RunMigrations {
		if c.Environment != "TEST" && c.Environment != "LOCAL" {
			if c.KafkaPrimaryCACert != "" && c.KafkaPrimaryPublicKey != "" && c.KafkaPrimaryPrivateKey != "" {
				c.KafkaPrimaryTLSConf, err = c.TLSConfigPrimary()
				if err != nil {
					return err
				}
			} else {
				c.KafkaPrimaryTLSConf = &tls.Config{MinVersion: tls.VersionTLS12}
			}

			if c.KafkaSecondaryCACert != "" && c.KafkaSecondaryPublicKey != "" && c.KafkaSecondaryPrivateKey != "" {
				c.KafkaSecondaryTLSConf, err = c.TLSConfigSecondary()
				if err != nil {
					return err
				}
			} else {
				c.KafkaSecondaryTLSConf = &tls.Config{MinVersion: tls.VersionTLS12}
			}

			if c.TemporalTLSConf, err = c.TemporalConfig.TLSConfig(); err != nil {
				return err
			}
		}
	}

	return err
}
