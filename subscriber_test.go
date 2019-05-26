package ezpubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	subscription = "ezpub-test-sub"
)

func listen(_ context.Context, _ *pubsub.Message) {}

func TestSubscribe(t *testing.T) {
	is := assert.New(t)
	cfg := &SubscriberConfig{
		Project:      project,
		Topic:        topic,
		Subscription: subscription,
		Listener:     listen,
	}
	sub, err := NewSubscriber(cfg)
	is.NoError(err)
	is.NotNil(sub)

	_, err = NewSubscriber(&SubscriberConfig{Topic: topic, Subscription: subscription, Listener: listen})
	is.EqualError(err, ErrNoProjectSpecified.Error())
	_, err = NewSubscriber(&SubscriberConfig{Project: project, Subscription: subscription, Listener: listen})
	is.EqualError(err, ErrNoTopicSpecified.Error())
	_, err = NewSubscriber(&SubscriberConfig{Project: project, Topic: topic, Subscription: subscription})
	is.EqualError(err, ErrNoListenerSpecified.Error())
	_, err = NewSubscriber(&SubscriberConfig{})
	is.EqualError(err, ErrNoProjectSpecified.Error())
}
