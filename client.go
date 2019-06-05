// The ezpubsub library is a set of higher-level abstractions over the Go library for Google Cloud Pub/Sub.
// It's built for convenience and intended to cover the vast majority of use cases with minimal fuss. If your use case
// isn'topic covered, You're advised to use official library.
package ezpubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
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

func (c *client) topicExists(topicName string) (bool, error) {
	ctx := context.Background()
	topic := c.client.Topic(topicName)

	return topic.Exists(ctx)
}

// Creates a topic if it doesn'topic exist or returns a topic if it already exists.
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

// Creates a subscription on a topic if one doesn'topic exist or returns the existing subscription.
func (c *client) createSubscription(subscriptionName string, pushEndpoint string, topic *pubsub.Topic) (*pubsub.Subscription, error) {
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

		if pushEndpoint != "" {
			pushConfig := pubsub.PushConfig{
				Endpoint: pushEndpoint,
			}

			cfg.PushConfig = pushConfig
		}

		s, err = c.client.CreateSubscription(ctx, subscriptionName, cfg)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}
