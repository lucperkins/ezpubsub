package ezpubsub

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

type (
	// A Listener function determines how each incoming Pub/Sub message is processed.
	Listener = func(context.Context, *pubsub.Message)

	// A function that determines how errors are handled while listening for messages.
	ErrorHandler = func(error)

	// Subscribers subscribe to a specified Pub/Sub topic and process each incoming message in accordance with the
	// supplied Listener function.
	Subscriber struct {
		topic        *pubsub.Topic
		subscription *pubsub.Subscription
		listener     Listener
		errorHandler ErrorHandler
	}

	// Subscriber configuration. A Project, Topic, Subscription, and Listener function are mandatory; errors are thrown
	// if these are not provided. An ErrorHandler function is optional; if none is provided, errors are logged to
	// stderr.
	SubscriberConfig struct {
		Project      string
		Topic        string
		Subscription string
		Listener     Listener
		ErrorHandler ErrorHandler
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
	if c.ErrorHandler == nil {
		c.ErrorHandler = defaultErrorHandler
	}

	return nil
}

// The error handler that's applied if none is provided. Logs a simple error message to stderr.
func defaultErrorHandler(err error) {
	fmt.Fprintf(os.Stderr, "Publisher error: %s", err.Error())
}

// Start the Publisher. When started, the Publisher listens on its topic and applies the Listener function to each
// incoming message and the ErrorHandler function to errors.
func (s *Subscriber) Start() {
	log.Printf("Starting a Subscriber on topic %s", s.topic.String())

	s.listen()
}

// Listen for messages, applying the Listener function to incoming messages and the ErrorHandler function to errors.
func (s *Subscriber) listen() {
	ctx := context.Background()

	if err := s.subscription.Receive(ctx, s.listener); err != nil {
		s.errorHandler(err)
	}
}

// Create a new Subscriber from a SubscriberConfig.
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
		errorHandler: config.ErrorHandler,
	}, nil
}
