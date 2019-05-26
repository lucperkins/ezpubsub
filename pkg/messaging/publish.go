package messaging

import (
	"cloud.google.com/go/pubsub"
	"context"
	"log"
)

type Publisher struct {
	topic    *pubsub.Topic
	receiver chan *pubsub.PublishResult
}

func NewPublisher(projectName, topicName string) (*Publisher, error) {
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

	return &Publisher{
		topic:    topic,
		receiver: rcv,
	}, nil
}

func (p *Publisher) Start() {
	for {
		res := <-p.receiver
		s, err := res.Get(context.Background())
		if err != nil {
			log.Printf("Error: %s", err.Error())
		}
		log.Printf("Server ID: %s", s)
	}
}
