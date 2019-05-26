package ezpubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"log"
)

type (
	// A Notifier function determines how message publishing results are processed.
	Notifier = func(*pubsub.PublishResult)

	// Publishers publish messages on a specified Pub/Sub topic.
	Publisher struct {
		topic    *pubsub.Topic
		notifier Notifier
	}

	// Publisher configuration.
	PublisherConfig struct {
		Project  string
		Topic    string
		Notifier Notifier
	}
)

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
	ctx := context.Background()

	err := config.validate()
	if err != nil {
		return nil, err
	}

	client, err := newClient(config.Project)
	if err != nil {
		return nil, err
	}

	topic, err := client.createTopic(ctx, config.Topic)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		topic:    topic,
		notifier: config.Notifier,
	}, nil
}

// Publish the specified data payload (as raw bytes) on the Publisher's topic.
func (p *Publisher) Publish(data []byte) {
	ctx := context.Background()

	log.Printf("Publishing a message to topic %s", p.topic.String())

	msg := &pubsub.Message{
		Data: data,
	}
	res := p.topic.Publish(ctx, msg)

	if p.notifier != nil {
		p.notifier(res)
	}

	defer p.topic.Stop()
}