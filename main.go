package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"log"
)

const (
	project = "shear-dev"
	topic = "all-msgs"
)

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

type producer struct {
	project string
	topic *pubsub.Topic
	receiver chan *pubsub.PublishResult
}

func newProducer(project, topicName string) (*producer, error) {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, project)

	if err != nil {
		return nil, err
	}

	topic := client.Topic(topicName)

	ok, err := topic.Exists(ctx)

	if err != nil {
		return nil, err
	}

	if !ok {
		topic, err = client.CreateTopic(ctx, topicName)
		if err != nil {
			return nil, err
		}
	}

	rcv := make(chan *pubsub.PublishResult, 1)

	return &producer{
		project: project,
		topic: topic,
		receiver: rcv,
	}, nil
}

func (p *producer) publish(ctx context.Context, data []byte) {
	msg := &pubsub.Message{
		Data: data,
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

func main() {
	producer, err := newProducer(project, topic)

	go func() {
		producer.start()
	}()

	panicOnErr(err)
	bs := []byte("Hello world")
	producer.publish(context.Background(), bs)
}
