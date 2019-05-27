package ezpubsub

import (
	"cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

const (
	subscription = "ezpub-test-sub"
)

func TestSubscribe(t *testing.T) {
	is := assert.New(t)
	cfg := &SubscriberConfig{
		Project:      project,
		Topic:        topic,
		Subscription: subscription,
	}
	sub, err := NewSubscriber(cfg)
	is.NoError(err)
	is.NotNil(sub)
	is.NotNil(sub.listener)
	is.NotNil(sub.errorHandler)
	is.NotNil(sub.topic)
	is.NotNil(sub.subscription)

	_, err = NewSubscriber(&SubscriberConfig{Topic: topic, Subscription: subscription})
	is.EqualError(err, ErrNoProjectSpecified.Error())
	_, err = NewSubscriber(&SubscriberConfig{Project: project, Subscription: subscription})
	is.EqualError(err, ErrNoTopicSpecified.Error())
	_, err = NewSubscriber(&SubscriberConfig{Project: project, Topic: topic})
	is.EqualError(err, ErrNoSubscriptionSpecified.Error())
	_, err = NewSubscriber(&SubscriberConfig{})
	is.EqualError(err, ErrNoProjectSpecified.Error())
}

func ExampleSubscriber() {
	subscriberConfig := &SubscriberConfig{
		Project:      "...",
		Topic:        "...",
		Subscription: "...",
		Listener: func(msg *pubsub.Message) {
			log.Printf("Message received (id: %s, payload: %s)", msg.Data, string(msg.Data))

			msg.Ack()
		},
		ErrorHandler: func(err error) {
			log.Printf("Publisher error: %v", err)
		},
	}

	subscriber, err := NewSubscriber(subscriberConfig)
	if err != nil {
		log.Fatalf("Subscriber creation error: %s", err)
	}

	subscriber.Start()
}
