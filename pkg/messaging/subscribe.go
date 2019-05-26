package messaging

import (
	"cloud.google.com/go/pubsub"
	"context"
	"log"
)

type listener = func(context.Context, *pubsub.Message)

type Subscriber struct {
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
	listener     listener
}

func (s *Subscriber) Start() {
	log.Printf("Starting a Subscriber on topic %s", s.topic.String())

	ctx := context.Background()
	err := s.subscription.Receive(ctx, s.listener)
	if err != nil {
		panic(err)
	}
}

func NewSubscriber(projectName, topicName, subscriptionName string, listener listener) (*Subscriber, error) {
	ctx := context.Background()

	client, err := newClient(projectName)
	if err != nil {
		return nil, err
	}

	topic, err := client.createTopic(ctx, topicName)
	if err != nil {
		return nil, err
	}

	sub, err := client.createSubscription(ctx, subscriptionName, topic)
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		topic:        topic,
		subscription: sub,
		listener: listener,
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, data []byte) {
	log.Printf("Publishing a message to topic %s", p.topic.String())

	msg := &pubsub.Message{
		Data: data,
	}
	res := p.topic.Publish(ctx, msg)

	p.receiver <- res

	defer p.topic.Stop()
}

