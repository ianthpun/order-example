package services

import (
	"os"
	"time"

	"github.com/axiomzen/envconfig"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type BaseConfig struct {
	LogLevel         string  `default:"INFO"`
	LogFormat        string  `default:"TEXT"`
	ServiceName      string  `required:"true"`
	Revision         string  `envconfig:"REVISION"`
	Environment      string  `default:"LOCAL" required:"true"`
	PodName          string  `required:"true"`
	SentryDSN        string  `envconfig:"SENTRY_DSN"`
	SentrySampleRate float64 `default:"1.0" envconfig:"SENTRY_SAMPLE_RATE"`
}

// EventLoopConfig configures eventloop behavior
type EventLoopConfig struct {
	EventLoopEnabled    bool          `default:"true"`
	EventHandlerTimeout time.Duration `default:"30s"`
}

type OAuthConfig struct {
	// OAuth config
	// Much of the config below is already specified in the files served from the following endpoints
	// Staging: https://dapperlabs.dapperlabsdev-private.auth0.com/.well-known/openid-configuration
	// Prod: https://dapperlabs.auth0.com/.well-known/openid-configuration

	// OAuthJWKSURL is the URL of the OAuth providers JSON web key set
	OAuthJWKSURL string `default:"https://dapperlabs.dapperlabsdev-private.auth0.com/.well-known/jwks.json"`
	// OAuthTokenAudience is the audience expected in access tokens.
	OAuthTokenAudience string `default:"https://api.staging.app.dapperlabs.com"`
	// OAuthTokenIssuer is the URL of the service that issues ID tokens to OAuth clients
	OAuthTokenIssuer string `default:"https://auth.staging.meetdapper.com/"`
	// AccessTokenHeader is the HTTP request header in which clients pass their ID tokens
	AccessTokenHeader string `default:"Authorization"`
	// MockAccessTokenHeader is the HTTP request header in which clients pass their mock ID tokens
	MockAccessTokenHeader string `default:"x-mock-access-token"`
}

// GRPCServiceConfig defines the config for grpc services
type GRPCServiceConfig struct {
	LivenessTolerancePeriod time.Duration `default:"60s"`
	GRPCEnabled             bool          `default:"true"`
	GRPCPort                uint16        `default:"8200"`
	GRPCHealthPort          uint16        `default:"8300"`
	GRPCMaxConnectionAge    time.Duration `default:"2m"`
}

type DependencyComputableConfig interface {
	ComputeDependents() error
}

// LoadConfig will attempt to load environment variables into the config struct passed in.
// If the config implements DependencyComputableConfig, ComputeDependents will be invoked.
func LoadConfig(logger *zap.Logger, prefix string, conf interface{}) error {
	initErr := envconfig.Process(prefix, conf)
	if initErr == nil {
		if computableConfig, ok := conf.(DependencyComputableConfig); ok {
			initErr = computableConfig.ComputeDependents()
		}
	}

	if initErr != nil {
		logger.Error("Error reading configuration", zap.Error(initErr))
	}
	return initErr
}

func LoadDotEnv(logger *zap.Logger, envKey string) {
	if envFileName := os.Getenv(envKey); envFileName != "" {
		if err := godotenv.Load(envFileName); err != nil {
			logger.Error("error loading godotenv", zap.Error(err))
		}
	}
}
