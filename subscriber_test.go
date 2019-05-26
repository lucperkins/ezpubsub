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

func receive(_ context.Context, msg *pubsub.Message) {}

func TestSubscribe(t *testing.T) {
	is := assert.New(t)
	sub, err := NewSubscriber(project, topic, subscription, receive)
	is.NoError(err)
	is.NotNil(sub)
}