package ezpubsub

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/thanhpk/randstr"
	"testing"
)

const (
	testProject = "any-project"
)

var (
	testTopic        = randstr.String(12)
	testSubscription = randstr.String(10)
)

func TestAdminInterface(t *testing.T) {
	is := assert.New(t)
	admin, err := NewAdmin(testProject)
	is.NoError(err)
	is.NotNil(admin)
	is.NotNil(admin.client)

	is.NoError(admin.DeleteSubscriptions(testSubscription))

	topics, err := admin.ListTopics()
	is.NoError(err)
	is.NotNil(topics)

	subscriptions, err := admin.ListSubscriptions()
	is.NoError(err)
	is.NotNil(subscriptions)
	is.Empty(subscriptions)

	exists, err := admin.TopicExists(testTopic)
	is.NoError(err)
	is.False(exists)

	cfg := &SubscriberConfig{
		Project:      testProject,
		Topic:        testTopic,
		Subscription: testSubscription,
	}
	sub, err := NewSubscriber(cfg)
	is.NoError(err)
	is.NotNil(sub)

	is.NoError(admin.DeleteSubscription(testSubscription))

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
