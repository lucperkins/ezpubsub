package ezpubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"log"
)

type (
	// A Listener function determines how each incoming Pub/Sub message is processed.
	Listener = func(context.Context, *pubsub.Message)

	// Subscribers subscribe to a specified Pub/Sub topic and process each incoming message in accordance with the
	// supplied listener function.
	Subscriber struct {
		topic        *pubsub.Topic
		subscription *pubsub.Subscription
		listener     Listener
	}

	// Subscriber configuration. All fields except `Listener` are mandatory.
	SubscriberConfig struct {
		Project      string
		Topic        string
		Subscription string
		Listener     Listener
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

// Start the Publisher. When started, the Publisher listens on its topic and applies its listener function to each
// incoming message.
func (s *Subscriber) Start() {
	log.Printf("Starting a Subscriber on topic %s", s.topic.String())

	ctx := context.Background()
	err := s.subscription.Receive(ctx, s.listener)
	if err != nil {
		panic(err)
	}
}

// Create a new Subscriber from a SubscriberConfig
func NewSubscriber(config *SubscriberConfig) (*Subscriber, error) {
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

	sub, err := client.createSubscription(config.Subscription, topic)
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		topic:        topic,
		subscription: sub,
		listener:     config.Listener,
	}, nil
}
