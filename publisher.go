package ezpubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
)

type (
	// A `Notifier` function determines how message publishing results are processed.
	Notifier = func(*pubsub.PublishResult)

	// `Publisher`s publish messages on a specified Pub/Sub topic.
	Publisher struct {
		topic    *pubsub.Topic
		notifier Notifier
	}

	// Publisher configuration. All fields except `Notifier` are mandatory.
	PublisherConfig struct {
		Project  string
		Topic    string
		Notifier Notifier
	}
)

// Validate the `PublisherConfig`
func (c *PublisherConfig) validate() error {
	if c.Project == "" {
		return ErrNoProjectSpecified
	}

	if c.Topic == "" {
		return ErrNoTopicSpecified
	}

	return nil
}

// Create a new `Publisher` from a `PublisherConfig`.
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
		topic:    topic,
		notifier: config.Notifier,
	}, nil
}

// Publish the specified `data` payload (as raw bytes) on the `Publisher`'s topic.
func (p *Publisher) Publish(data []byte) {
	ctx := context.Background()

	msg := &pubsub.Message{
		Data: data,
	}
	res := p.topic.Publish(ctx, msg)

	if p.notifier != nil {
		p.notifier(res)
	}

	defer p.topic.Stop()
}

// Publish a batch of messages synchronously.
func (p *Publisher) BatchPublishSync(payloads [][]byte) {
	msgs := convertDataToMessages(payloads)

	for _, m := range msgs {
		res := p.publishMessage(m)
		p.notify(res)
	}
}

// Publish a batch of messages asynchronously.
func (p *Publisher) BatchPublishAsync(payloads [][]byte) {
	ms := convertDataToMessages(payloads)

	done := make(chan bool, 1)

	go p.asyncWorker(ms, done)

	<-done
}

// Converts a slice of raw data payloads into a slice of Messages
func convertDataToMessages(payloads [][]byte) []*pubsub.Message {
	msgs := make([]*pubsub.Message, len(payloads))

	for _, p := range payloads {
		msg := &pubsub.Message{
			Data: p,
		}
		msgs = append(msgs, msg)
	}
	return msgs
}

// Publish a message on the Publisher's topic.
func (p *Publisher) publishMessage(msg *pubsub.Message) *pubsub.PublishResult {
	ctx := context.Background()
	return p.topic.Publish(ctx, msg)
}

// Apply the notification function if one is specified.
func (p *Publisher) notify(res *pubsub.PublishResult) {
	if p.notifier != nil {
		p.notifier(res)
	}
}

// If the number processed equals the total number of messages in the batch, notify the channel that the processing is
// done.
func notifyWhenDone(published, queue int, done chan bool) {
	if published == queue {
		done <- true
	}
}

// A channel-based async worker for batch message publishing.
func (p *Publisher) asyncWorker(messages []*pubsub.Message, done chan bool) {
	queueLength := len(messages)
	numPublished := 0

	for _, m := range messages {
		res := p.publishMessage(m)
		p.notify(res)

		numPublished += 1

		notifyWhenDone(numPublished, queueLength, done)
	}
}
