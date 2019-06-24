package ezpubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"google.golang.org/api/iterator"
)

type (
	// A simple administrative interface for Pub/Sub projects.
	Admin struct {
		client *client
	}
)

// Create a new Admin, specifying the Google Cloud project name.
func NewAdmin(project string) (*Admin, error) {
	client, err := newClient(project)
	if err != nil {
		return nil, err
	}

	return &Admin{
		client: client,
	}, nil
}

// List all current topics under the specified project.
func (a *Admin) ListTopics() ([]string, error) {
	ctx := context.Background()

	it := a.client.client.Topics(ctx)

	return topicIteratorToList(it)
}

// Checks is a topic already exists.
func (a *Admin) TopicExists(topicName string) (bool, error) {
	return a.client.topicExists(topicName)
}

// Deletes the specified topic.
func (a *Admin) DeleteTopic(topicName string) error {
	ctx := context.Background()
	return a.client.client.Topic(topicName).Delete(ctx)
}

// Checks if a subscription exists.
func (a *Admin) SubscriptionExists(subscriptionName string) (bool, error) {
	ctx := context.Background()
	return a.client.client.Subscription(subscriptionName).Exists(ctx)
}

// Lists all current subscriptions.
func (a *Admin) ListSubscriptions() ([]string, error) {
	ctx := context.Background()

	it := a.client.client.Subscriptions(ctx)

	return subscriptionIteratorToList(it)
}

// Deletes a specified subscription.
func (a *Admin) DeleteSubscription(subscription string) error {
	ctx := context.Background()
	return a.client.client.Subscription(subscription).Delete(ctx)
}

// Deletes multiple subscriptions.
func (a *Admin) DeleteSubscriptions(subscriptions ...string) error {
	for _, sub := range subscriptions {
		if err := a.DeleteSubscription(sub); err != nil {
			return err
		}
	}

	return nil
}

func topicIteratorToList(it *pubsub.TopicIterator) ([]string, error) {
	ts := make([]string, 0)

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

func subscriptionIteratorToList(it *pubsub.SubscriptionIterator) ([]string, error) {
	ts := make([]string, 0)

	for {
		sub, err := it.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		ts = append(ts, sub.String())
	}

	return ts, nil
}
