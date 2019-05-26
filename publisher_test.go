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
}
