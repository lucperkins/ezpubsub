package messaging

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	subscription = "my-sub"
)

func listen(_ context.Context, _ *pubsub.Message) {}

func TestSubscribe(t *testing.T) {
	is := assert.New(t)
	cfg := &SubscriberConfig{
		Project: project,
		Topic: topic,
		Subscription: subscription,
		Listener: listen,
	}
	sub, err := NewSubscriber(cfg)
	is.NoError(err)
	is.NotNil(sub)
}