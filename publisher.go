package ezpubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
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

// Synchronously publish a batch of message payloads, preserving message order.
func (p *Publisher) PublishBatchSync(payloads [][]byte) {
	ctx := context.Background()
	msgs := convertDataToMessages(payloads)

	for _, msg := range msgs {
		if p.n != nil {
			res := p.t.Publish(ctx, msg)
			p.n(res)
		} else {
			p.t.Publish(ctx, msg)
		}
	}
}

// Converts a slice of raw data payloads into a slice of Messages
func convertDataToMessages(payloads [][]byte) []*pubsub.Message {
	msgs := make([]*pubsub.Message, 0)

	for _, p := range payloads {
		msg := &pubsub.Message{
			Data: p,
		}
		msgs = append(msgs, msg)
	}
	return msgs
}
