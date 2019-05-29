// The ezpubsub library is a set of higher-level abstractions over the Go library for Google Cloud Pub/Sub.
// It's built for convenience and intended to cover the vast majority of use cases with minimal fuss. If your use case
// isn't covered, You're advised to use official library.
package ezpubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"google.golang.org/api/iterator"
)

type client struct {
	client *pubsub.Client
}

// Creates a new Pub/Sub client or throws an error.
func newClient(project string) (*client, error) {
	ctx := context.Background()
	cl, err := pubsub.NewClient(ctx, project)

	if err != nil {
		return nil, err
	}

	return &client{
		client: cl,
	}, nil
}

// Creates a topic if it doesn't exist or returns a topic if it already exists.
func (c *client) createTopic(topicName string) (*pubsub.Topic, error) {
	var topic *pubsub.Topic
	ctx := context.Background()

	exists, err := c.topicExists(topicName)
	if err != nil {
		return nil, err
	}

	if !exists {
		topic, err = c.client.CreateTopic(ctx, topicName)
		if err != nil {
			return nil, err
		}
	} else {
		topic = c.client.Topic(topicName)
	}

	return topic, nil
}

// Checks if the topic already exists.
func (c *client) topicExists(topicName string) (bool, error) {
	ctx := context.Background()
	topic := c.client.Topic(topicName)

	return topic.Exists(ctx)
}

// Lists all topics associated with a project.
func (c *client) listTopics() ([]string, error) {
	ctx := context.Background()
	ts := make([]string, 0)

	it := c.client.Topics(ctx)

	for {
		topic, err := it.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		ts = append(ts, topic.String())
	}

	return ts, nil
}

// Lists all current subscriptions
func (c *client) listSubscriptions() ([]string, error) {
	ctx := context.Background()
	ss := make([]string, 0)

	it := c.client.Subscriptions(ctx)

	for {
		sub, err := it.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		ss = append(ss, sub.String())
	}

	return ss, nil
}


// Creates a subscription on a topic if one doesn't exist or returns the existing subscription.
func (c *client) createSubscription(subscriptionName string, topic *pubsub.Topic) (*pubsub.Subscription, error) {
	ctx := context.Background()

	s := c.client.Subscription(subscriptionName)
	exists, err := s.Exists(ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
		cfg := pubsub.SubscriptionConfig{
			Topic: topic,
		}

		s, err = c.client.CreateSubscription(ctx, subscriptionName, cfg)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}
