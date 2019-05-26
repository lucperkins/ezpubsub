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
	pub, err := NewPublisher(project, topic)
	is.NoError(err)
	is.NotNil(pub)
}
