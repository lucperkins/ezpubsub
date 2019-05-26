package messaging

import (
	"cloud.google.com/go/pubsub"
	"context"
	"log"
)

type (
	listenerFunc = func(ctx context.Context, msg *pubsub.Message)

	Subscriber struct {
		topic        *pubsub.Topic
		subscription *pubsub.Subscription
		listener     listenerFunc
	}

	SubscriberConfig struct {
		Project      string
		Topic        string
		Subscription string
		Listener     listenerFunc
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
	log.Printf("Starting a subscriber on topic %s", s.topic.String())

	ctx := context.Background()
	err := s.subscription.Receive(ctx, s.listener)
	if err != nil {
		panic(err)
	}
}

func NewSubscriber(config *SubscriberConfig) (*Subscriber, error) {
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

	sub, err := client.createSubscription(ctx, config.Subscription, topic)
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		topic:        topic,
		subscription: sub,
		listener:     config.Listener,
	}, nil
}
