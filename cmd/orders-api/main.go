package main

import (
	eventclient2 "order-sample/pkg/eventclient"
	kafka2 "order-sample/pkg/eventclient/kafka"
)

type handler struct{}

func (h *handler) handleEvent(msg *eventclient2.Message) ([]*eventclient2.Message, error) {
	return nil, nil
}

func main() {
	router := eventclient2.NewRouter()

	publisher, err := kafka2.NewPublisher(
		[]string{"broker"},
	)
	if err != nil {
		panic(err)
	}

	subscriber, err := kafka2.NewSubscriber(
		[]string{"broker"},
		"consumer_group_test",
	)
	if err != nil {
		panic(err)
	}

	router.Register(
		"test_handler",
		"SOME_TOPIC",
		subscriber,
		"SOME_TOPIC_TO_SEND",
		publisher,
		handler{}.handleEvent,
	)
}
