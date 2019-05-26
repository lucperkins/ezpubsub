package ezpubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"log"
)

type (
	notifierFunc = func(*pubsub.PublishResult)

	publisher struct {
		topic    *pubsub.Topic
		notifier notifierFunc
	}

	// Publisher configuration
	PublisherConfig struct {
		Project  string
		Topic    string
		Notifier notifierFunc
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

// Create a new publisher from a PublisherConfig
func NewPublisher(config *PublisherConfig) (*publisher, error) {
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

	return &publisher{
		topic:    topic,
		notifier: config.Notifier,
	}, nil
}

// Publish a data payload (as raw bytes) on the Publisher's topic
func (p *publisher) Publish(data []byte) {
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