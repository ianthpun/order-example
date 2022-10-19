package domain

type Message struct {
	id      string
	topic   string
	payload []byte
}

func NewMessage(
	id string,
	topic string,
	payload []byte,
) Message {
	return Message{
		id:      id,
		topic:   topic,
		payload: payload,
	}
}

func (m *Message) GetID() string {
	return m.id
}

func (m *Message) GetTopic() string {
	return m.topic
}

func (m *Message) GetPayload() []byte {
	return m.payload
}
