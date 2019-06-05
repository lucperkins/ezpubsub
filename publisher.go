package ezpubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
)

type (
	// A handler for the server ID returned when publishing a message.
	ServerIdHandler = func(string)

	// Publishers publish messages on a specified Pub/Sub topic.
	Publisher struct {
		topic           *pubsub.Topic
		errorHandler    ErrorHandler
		serverIdHandler ServerIdHandler
	}
)

// Create a new Publisher from a PublisherConfig.
func NewPublisher(config *PublisherConfig) (*Publisher, error) {
	err := config.validate()
	if err != nil {
		return nil, err
	}

	client, err := newClient(config.Project)
	if err != nil {
		return nil, err
	}

	topic, err := client.createTopic(config.Topic)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		topic:           topic,
		errorHandler:    config.ErrorHandler,
		serverIdHandler: config.ServerIdHandler,
	}, nil
}

// Publish the specified payload on the Publisher's topic.
func (p *Publisher) Publish(data []byte) {
	ctx := context.Background()

	msg := &pubsub.Message{
		Data: data,
	}

	res := p.topic.Publish(ctx, msg)

	p.handleResponse(ctx, res)
}

func (p *Publisher) handleResponse(ctx context.Context, res *pubsub.PublishResult) {
	id, err := res.Get(ctx)
	if err != nil && p.errorHandler != nil {
		p.errorHandler(err)
	}

	if p.serverIdHandler != nil {
		p.serverIdHandler(id)
	}

}

// Publish a JSON-serializable object on the Publisher's topic and throw an error if JSON marshalling is unsuccessful.
func (p *Publisher) PublishObject(obj interface{}) error {
	bs, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	p.Publish(bs)
	return nil
}

// Publish a string on the Publisher's topic.
func (p *Publisher) PublishString(s string) {
	p.Publish([]byte(s))
}

// Synchronously publish a batch of message payloads, preserving message order.
func (p *Publisher) PublishBatchSync(payloads [][]byte) {
	msgs := convertDataToMessages(payloads)

	for _, msg := range msgs {
		p.Publish(msg.Data)
	}
}

// Converts a slice of raw data payloads into a slice of Messages
func convertDataToMessages(payloads [][]byte) []*pubsub.Message {
	msgs := make([]*pubsub.Message, 0)

	for _, p := range payloads {
		msg := &pubsub.Message{
			Data: p,
		}
		msgs = append(msgs, msg)
	}
	return msgs
}
