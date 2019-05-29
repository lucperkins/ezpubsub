package ezpubsub

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
	"testing"
)

func TestAdminInterface(t *testing.T) {
	is := assert.New(t)
	admin, err := NewAdmin("any-project")
	is.NoError(err)
	is.NotNil(admin)
	is.NotNil(admin.client)

	topics, err := admin.ListTopics()
	is.NoError(err)
	is.NotNil(topics)

	subscriptions, err := admin.ListSubscriptions()
	is.NoError(err)
	is.NotNil(subscriptions)

	randTopic := randstr.String(10)
	exists, err := admin.TopicExists(randTopic)
	is.NoError(err)
	is.False(exists)
}

func ExampleAdmin() {
	admin, err := NewAdmin("my-project")
	if err != nil {
		// handle error
	}

	topics, err := admin.ListTopics()
	if err != nil {
		// handle error
	}

	fmt.Println("Listing topics:")
	for _, topic := range topics {
		fmt.Println(topic)
	}
}
