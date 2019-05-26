package ezpubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"log"
)

type (
	listenerFunc = func(ctx context.Context, msg *pubsub.Message)

	subscriber struct {
		topic        *pubsub.Topic
		subscription *pubsub.Subscription
		listener     listenerFunc
	}

	// Subscriber configuration
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

// Start the publisher. When started, the publisher listens on its topic and applies its listener function to each
// incoming message.
func (s *subscriber) Start() {
	log.Printf("Starting a subscriber on topic %s", s.topic.String())

	ctx := context.Background()
	err := s.subscription.Receive(ctx, s.listener)
	if err != nil {
		panic(err)
	}
}

// Create a new Subscriber from a SubscriberConfig
func NewSubscriber(config *SubscriberConfig) (*subscriber, error) {
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

	return &subscriber{
		topic:        topic,
		subscription: sub,
		listener:     config.Listener,
	}, nil
}
