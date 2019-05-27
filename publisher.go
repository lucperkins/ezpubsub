package ezpubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
)

type (
	// A function that specifies what happens when a message is published.
	Notifier = func(*pubsub.PublishResult)

	// Publishers publish messages on a specified Pub/Sub t.
	Publisher struct {
		t *pubsub.Topic
		n Notifier
	}

	// Publisher configuration. All fields except Notifier are mandatory.
	PublisherConfig struct {
		Project  string
		Topic    string
		Notifier Notifier
	}
)

// Validate the PublisherConfig
func (c *PublisherConfig) validate() error {
	if c.Project == "" {
		return ErrNoProjectSpecified
	}

	if c.Topic == "" {
		return ErrNoTopicSpecified
	}

	return nil
}

// Create a new Publisher from a PublisherConfig.
func NewPublisher(config *PublisherConfig) (*Publisher, error) {
	err := config.validate()
	if err != nil {
		return nil, err
	}

	client, err := newClient(config.Project)
	if err != nil {
		return nil, err
	}

	topic, err := client.createTopic(config.Topic)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		t: topic,
		n: config.Notifier,
	}, nil
}

// Publish the specified data payload (as raw bytes) on the Publisher's topic.
func (p *Publisher) Publish(data []byte) {
	ctx := context.Background()
	defer p.t.Stop()

	msg := &pubsub.Message{
		Data: data,
	}

	if p.n != nil {
		res := p.t.Publish(ctx, msg)
		p.n(res)
	} else {
		p.t.Publish(ctx, msg)
	}
}

func (p *Publisher) PublishBatchSync(payloads [][]byte) {
	ctx := context.Background()
	msgs := convertDataToMessages(payloads)

	for _, msg := range msgs {
		if p.n != nil {
			res := p.t.Publish(ctx, msg)
			p.n(res)
			p.t.Stop()
		} else {
			p.t.Publish(ctx, msg)
			p.t.Stop()
		}
	}
}

// Converts a slice of raw data payloads into a slice of Messages
func convertDataToMessages(payloads [][]byte) []*pubsub.Message {
	msgs := make([]*pubsub.Message, len(payloads))

	for _, p := range payloads {
		msg := &pubsub.Message{
			Data: p,
		}
		msgs = append(msgs, msg)
	}
	return msgs
}

// If the number processed equals the total number of messages in the batch, notify the channel that the processing is
// done.
func notifyWhenDone(published, queue int, done chan bool) {
	if published == queue {
		done <- true
	}
}

// A channel-based async worker for batch message publishing.
func (p *Publisher) asyncWorker(messages []*pubsub.Message, done chan bool) {
	ctx := context.Background()
	queueLength := len(messages)

	fmt.Printf("Queue length: %d", queueLength)
	numPublished := 0

	for _, m := range messages {
		numPublished += 1

		fmt.Printf("Num published: %d", numPublished)

		p.t.Publish(ctx, m)

		notifyWhenDone(numPublished, queueLength, done)
	}
}
