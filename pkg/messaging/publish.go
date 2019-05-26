package messaging

import (
	"cloud.google.com/go/pubsub"
	"context"
	"errors"
	"log"
)

type (
	notifier = func(*pubsub.PublishResult)

	publisher struct {
		t *pubsub.Topic
		n notifier
	}

	PublisherConfig struct {
		Project string
		Topic string
		Notifier notifier
	}
)

func (c *PublisherConfig) validate() error {
	if c.Project == "" {
		return errors.New("no project specified")
	}

	if c.Topic == "" {
		return errors.New("no t specified")
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
		t: topic,
		n: config.Notifier,
	}, nil
}

func (p *publisher) Publish(ctx context.Context, data []byte) {
	log.Printf("Publishing a message to t %s", p.t.String())

	msg := &pubsub.Message{
		Data: data,
	}
	res := p.t.Publish(ctx, msg)

	if p.n != nil {
		p.n(res)
	}

	defer p.t.Stop()
}