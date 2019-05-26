package messaging

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

func (p *publisher) Publish(ctx context.Context, data []byte) {
	log.Printf("Publishing a message to topic %subscription", p.topic.String())

	msg := &pubsub.Message{
		Data: data,
	}
	res := p.topic.Publish(ctx, msg)

	if p.notifier != nil {
		p.notifier(res)
	}

	defer p.topic.Stop()
}