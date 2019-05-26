package messaging

import (
	"cloud.google.com/go/pubsub"
	"context"
)

type client struct {
	client *pubsub.Client
}

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

func (c *client) createTopic(ctx context.Context, topicName string) (*pubsub.Topic, error) {
	topic := c.client.Topic(topicName)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}

	if !exists {
		topic, err = c.client.CreateTopic(ctx, topicName)
		if err != nil {
			return nil, err
		}
	}

	return topic, nil
}

func (c *client) createSubscription(ctx context.Context, subscriptionName string, topic *pubsub.Topic) (*pubsub.Subscription, error) {
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
