package main

import (
	"context"
	"github.com/dapperlabs/dibs/v2/eventclient"
	"go.uber.org/zap"
	"order-sample/cmd/orders-api/internal/config"
	"order-sample/cmd/orders-api/internal/ports/temporal/eventloops"
	"order-sample/internal/services"
)

func main() {
	ctx := context.Background()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	// load config from environment
	var conf config.Config
	if err := services.LoadConfig(logger, "ORDERS", &conf); err != nil {
		panic(err)
	}

	if conf.EventLoopEnabled {
		eventloops.New(eventloops.Config{
			Logger:           logger,
			KafkaClient:      nil,
			EventClientConf:  eventclient.Config{},
			WorkflowExecutor: nil,
		}).Start(ctx)
	}

}
