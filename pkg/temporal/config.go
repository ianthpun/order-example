package temporal

import (
	"crypto/tls"
	"fmt"
)

// Config is the Temporal configuration
type Config struct {
	TemporalHostAddr   string `default:"temporal:7233"`
	TemporalNamespace  string `default:"default"`
	TemporalServerName string
	TemporalPrivateKey string
	TemporalPublicKey  string

	// Computed config.
	// It's important that `tls.Config` is a pointer value because it contains a `Mutex` that cannot be safely passed
	// by value.
	TemporalTLSConf *tls.Config `ignored:"true"`

	// Temporal worker enable config
	TemporalWorkerEnabled bool `default:"true"`
}

// TLSConfig returns the TLS configuration
func (c *Config) TLSConfig() (*tls.Config, error) {
	cert, err := tls.X509KeyPair([]byte(c.TemporalPublicKey), []byte(c.TemporalPrivateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse Temporal credentials: %w", err)
	}

	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
		ServerName:   c.TemporalServerName,
	}, nil
}
