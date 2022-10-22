package watermill

import (
	"github.com/ThreeDotsLabs/watermill/message"
	eventclient2 "order-sample/pkg/eventclient"
)

func ToWatermillMiddlewareFunc(mid eventclient2.MiddlewareFunc) message.HandlerMiddleware {
	return func(handlerFunc message.HandlerFunc) message.HandlerFunc {
		return handlerFunc
	}
}

func ToWatermillHandlerFunc(f eventclient2.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		msgs, err := f(WatermillMessageToMessage(msg))
		if err != nil {
			return nil, err
		}

		var watermillMsgs []*message.Message

		for _, msg := range msgs {
			watermillMsgs = append(watermillMsgs, MessageToWatermillMessage(msg))
		}

		return watermillMsgs, nil
	}
}

func WatermillMessageToMessage(m *message.Message) *eventclient2.Message {
	if m == nil {
		return nil
	}

	return &eventclient2.Message{
		ID:           m.UUID,
		PartitionKey: m.Metadata.Get("partition"),
		Payload:      m.Payload,
	}
}

func MessageToWatermillMessage(m *eventclient2.Message) *message.Message {
	if m == nil {
		return nil
	}

	metadata := make(message.Metadata)
	metadata.Set(eventclient2.PartitionKey, m.PartitionKey)

	return &message.Message{
		UUID:     m.ID,
		Metadata: metadata,
		Payload:  m.Payload,
	}
}

func MessagesToWatermillMessages(messages []*eventclient2.Message) []*message.Message {
	var msgs []*message.Message
	for _, m := range messages {
		msgs = append(msgs, MessageToWatermillMessage(m))
	}

	return msgs
}

func WatermillMessagesToMessages(messages []*message.Message) []*eventclient2.Message {
	var msgs []*eventclient2.Message
	for _, m := range messages {
		msgs = append(msgs, WatermillMessageToMessage(m))
	}

	return msgs
}
