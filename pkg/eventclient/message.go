package eventclient

type Message struct {
	ID           string
	PartitionKey string // used if the pubsub library can partition based on key
	Payload      []byte
}

type Messages []*Message

func (msgs Messages) IDs() (ids []string) {
	for _, msg := range msgs {
		ids = append(ids, msg.ID)
	}

	return
}

const (
	PartitionKey = "partition_key"
)

func NewMessage(id string, payload []byte) *Message {
	return &Message{
		ID:      id,
		Payload: payload,
	}
}

func NewMessageWithPartitionKey(id string, partitionKey string, payload []byte) *Message {
	return &Message{
		ID:           id,
		PartitionKey: partitionKey,
		Payload:      payload,
	}
}
