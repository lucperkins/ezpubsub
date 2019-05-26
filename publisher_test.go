package ezpubsub

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

const (
	project = "ezpub"
	topic   = "ezpub-test-t"
)

func TestPublisher(t *testing.T) {
	is := assert.New(t)
	cfg := &PublisherConfig{
		Project: project,
		Topic:   topic,
	}
	pub, err := NewPublisher(cfg)
	is.NoError(err)
	is.NotNil(pub)

	_, err = NewPublisher(&PublisherConfig{Topic: topic})
	is.EqualError(err, ErrNoProjectSpecified.Error())
	_, err = NewPublisher(&PublisherConfig{Project: project})
	is.EqualError(err, ErrNoTopicSpecified.Error())
}

func ExamplePublisher() {
	publisherConfig := &PublisherConfig{
		Project: "...",
		Topic:   "...",
	}

	publisher, err := NewPublisher(publisherConfig)
	if err != nil {
		log.Fatalf("Publisher creation error: %s", err)
	}

	msg := []byte("Hello world")
	publisher.Publish(msg)
}
