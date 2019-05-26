package messaging

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	project = "shear-dev"
	topic = "shear-dev-test"
)

func TestPublisher(t *testing.T) {
	is := assert.New(t)
	cfg := &PublisherConfig{
		Project: project,
		Topic: topic,
	}
	pub, err := NewPublisher(cfg)
	is.NoError(err)
	is.NotNil(pub)

	_, err = NewPublisher(&PublisherConfig{Topic: topic})
	is.EqualError(err, ErrNoProjectSpecified.Error())
	_, err = NewPublisher(&PublisherConfig{Project: project})
	is.EqualError(err, ErrNoTopicSpecified.Error())
}
