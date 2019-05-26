package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"flag"
	"fmt"
	"log"
)

type publishable interface {
	bytes() []byte
	publish(*pubsub.Topic) *pubsub.PublishResult
}

type person struct {
	name string
	age int
}

func (p *person) bytes() []byte {
	return []byte(fmt.Sprintf("name: %s. age: %d.", p.name, p.age))
}

func (p *person) publish(t *pubsub.Topic) *pubsub.PublishResult {
	msg := &pubsub.Message{
		Data: p.bytes(),
	}
	return t.Publish(context.Background(), msg)
}

var _ publishable = (*person)(nil)

const (
	project = "shear-dev"
	topic = "all-msgs"
	subscription = "my-sub"
)

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

type producer struct {
	topic *pubsub.Topic
	receiver chan *pubsub.PublishResult
}

type subscriber struct {
	topic *pubsub.Topic
	subscription *pubsub.Subscription
}

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

func (s *subscriber) start() {
	log.Printf("Starting a subscriber on topic %s", s.topic.String())

	ctx := context.Background()
	err := s.subscription.Receive(ctx, receive)
	if err != nil {
		panic(err)
	}
}

func receive(_ context.Context, m *pubsub.Message) {
	s := string(m.Data)
	fmt.Printf("Message received: %s\n", s)
	m.Ack()
}

func newSubscriber(projectName, topicName, subscriptionName string) (*subscriber, error) {
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

	return &subscriber{
		topic: topic,
		subscription: sub,
	}, nil
}

func newProducer(projectName, topicName string) (*producer, error) {
	ctx := context.Background()

	client, err := newClient(projectName)
	if err != nil {
		return nil, err
	}

	topic, err := client.createTopic(ctx, topicName)
	if err != nil {
		return nil, err
	}

	rcv := make(chan *pubsub.PublishResult, 1)

	return &producer{
		topic: topic,
		receiver: rcv,
	}, nil
}

func (p *producer) publish(ctx context.Context, pub publishable) {
	msg := &pubsub.Message{
		Data: pub.bytes(),
	}
	res := p.topic.Publish(ctx, msg)

	p.receiver <- res

	defer p.topic.Stop()
}

func (p *producer) start() {
	for {
		res := <- p.receiver
		s, err := res.Get(context.Background())
		panicOnErr(err)
		log.Printf("Server ID: %s", s)
	}
}

var run string

func init() {
	flag.StringVar(&run, "run", "producer", "Which process type to run")
}

func main() {
	flag.Parse()

	if run == "producer" {
		log.Print("Running a producer")
		producer, err := newProducer(project, topic)
		panicOnErr(err)

		luc := person{
			name: "Luc",
			age: 37,
		}

		producer.publish(context.Background(), &luc)
	}

	if run == "subscriber" {
		log.Print("Running a subscriber")
		subscriber, err := newSubscriber(project, topic, subscription)
		panicOnErr(err)

		subscriber.start()
	}
}
