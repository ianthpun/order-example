package tests_test

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math/rand"
	eventclient2 "order-sample/pkg/eventclient"
	kafka2 "order-sample/pkg/eventclient/kafka"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var defaultTimeout = time.Second * 15

func init() {
	rand.Seed(3)
}

func kafkaBrokers() []string {
	return []string{"localhost:9091"}
}

func newPubSub(t *testing.T, consumerGroup string) (*kafka2.Publisher, *kafka2.Subscriber) {
	var err error
	var publisher *kafka2.Publisher

	retriesLeft := 5
	for {
		publisher, err = kafka2.NewPublisher(
			kafkaBrokers(),
		)
		if err == nil || retriesLeft == 0 {
			break
		}

		retriesLeft--
		fmt.Printf("cannot create kafka Publisher: %s, retrying (%d retries left)", err, retriesLeft)
		time.Sleep(time.Second * 2)
	}
	require.NoError(t, err)

	var subscriber *kafka2.Subscriber

	retriesLeft = 5
	for {
		subscriber, err = kafka2.NewSubscriber(
			kafkaBrokers(),
			consumerGroup,
		)
		if err == nil || retriesLeft == 0 {
			break
		}

		retriesLeft--
		fmt.Printf("cannot create kafka Subscriber: %s, retrying (%d retries left)", err, retriesLeft)
		time.Sleep(time.Second * 2)
	}

	require.NoError(t, err)

	return publisher, subscriber
}

func TestPublishSubscribe(
	t *testing.T,
) {
	pub, sub := newPubSub(t, "some_consumer_group")

	topicName := "test_topic"

	var messagesToPublish []*eventclient2.Message
	messagesPayloads := map[string][]byte{}

	for i := 0; i < 100; i++ {
		id := uuid.NewString()

		payload := []byte(fmt.Sprintf("%d", i))
		msg := eventclient2.NewMessage(id, payload)

		messagesToPublish = append(messagesToPublish, msg)
		messagesPayloads[id] = payload
	}
	err := publishWithRetry(pub, topicName, messagesToPublish...)
	require.NoError(t, err, "cannot publish message")

	messages, err := sub.Subscribe(context.Background(), topicName)
	require.NoError(t, err)

	receivedMessages, all := bulkRead("some message", messages, len(messagesToPublish), defaultTimeout*3)
	assert.True(t, all)

	AssertAllMessagesReceived(t, messagesToPublish, receivedMessages)
	//AssertMessagesPayloads(t, messagesPayloads, receivedMessages)
	//AssertMessagesMetadata(t, "test", messagesTestMetadata, receivedMessages)
	//
	//closePubSub(t, pub, sub)
	//assertMessagesChannelClosed(t, messages)
}

// AssertAllMessagesReceived checks if all messages were received,
// ignoring the order and assuming that they are already deduplicated.
func AssertAllMessagesReceived(t *testing.T, sent eventclient2.Messages, received eventclient2.Messages) bool {
	sentIDs := sent.IDs()
	receivedIDs := received.IDs()

	sort.Strings(sentIDs)
	sort.Strings(receivedIDs)

	assert.Equal(
		t,
		len(sentIDs), len(receivedIDs),
		"id's count is different: received: %d, sent: %d", len(receivedIDs), len(sentIDs),
	)

	return assert.Equal(
		t, sentIDs, receivedIDs,
		"received different messages ID's, missing: %s, extra %s",
		MissingMessages(sent, received),
		MissingMessages(received, sent),
	)
}
func difference(a, b []string) []string {
	mb := map[string]bool{}
	for _, x := range b {
		mb[x] = true
	}
	ab := []string{}
	for _, x := range a {
		if _, ok := mb[x]; !ok {
			ab = append(ab, x)
		}
	}
	return ab
}

// MissingMessages returns a list of missing messages UUIDs.
func MissingMessages(expected eventclient2.Messages, received eventclient2.Messages) []string {
	sentIDs := expected.IDs()
	receivedIDs := received.IDs()

	sort.Strings(sentIDs)
	sort.Strings(receivedIDs)

	return difference(sentIDs, receivedIDs)
}

func publishWithRetry(publisher eventclient2.Publisher, topic string, messages ...*eventclient2.Message) error {
	retries := 5

	for {
		err := publisher.Publish(topic, messages...)
		if err == nil {
			return nil
		}
		retries--

		fmt.Printf("error on publish: %s, %d retries left\n", err, retries)

		if retries == 0 {
			return err
		}
	}
}
func bulkRead(testID string, messagesCh <-chan *eventclient2.Message, limit int, timeout time.Duration) (receivedMessages eventclient2.Messages, all bool) {
	start := time.Now()

	defer func() {
		duration := time.Since(start)

		logMsg := "all messages (%d/%d) received in bulk read after %s of %s (test ID: %s)\n"
		if !all {
			logMsg = "not " + logMsg
		}

		log.Printf(logMsg, len(receivedMessages), limit, duration, timeout, testID)
	}()

MessagesLoop:
	for len(receivedMessages) < limit {
		select {
		case msg, ok := <-messagesCh:
			if !ok {
				break MessagesLoop
			}

			receivedMessages = append(receivedMessages, msg)
		case <-time.After(timeout):
			break MessagesLoop
		}
	}

	return receivedMessages, len(receivedMessages) == limit
}
