package messaging

import (
	"cloud.google.com/go/pubsub"
	"context"
	"log"
)

type (
	listener = func(ctx context.Context, msg *pubsub.Message)

	Subscriber struct {
		t *pubsub.Topic
		s *pubsub.Subscription
		l listener
	}

	SubscriberConfig struct {
		Project string
		Topic string
		Subscription string
		Listener listener
	}
)

func (c *SubscriberConfig) validate() error {
	if c.Project == "" {
		return ErrNoProjectSpecified
	}
	if c.Topic == "" {
		return ErrNoTopicSpecified
	}
	if c.Subscription == "" {
		return ErrNoSubscriptionSpecified
	}
	if c.Listener == nil {
		return ErrNoListenerSpecified
	}
	return nil
}

func (s *Subscriber) Start() {
	log.Printf("Starting a Subscriber on t %s", s.t.String())

	ctx := context.Background()
	err := s.s.Receive(ctx, s.l)
	if err != nil {
		panic(err)
	}
}

func NewSubscriber(config *SubscriberConfig) (*Subscriber, error) {
	ctx := context.Background()

	client, err := newClient(config.Project)
	if err != nil {
		return nil, err
	}

	topic, err := client.createTopic(ctx, config.Topic)
	if err != nil {
		return nil, err
	}

	sub, err := client.createSubscription(ctx, config.Subscription, topic)
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		t: topic,
		s: sub,
		l: config.Listener,
	}, nil
}
