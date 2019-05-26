package pub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/lucperkins/ezpubsub"
	"log"
)

func notify(res *pubsub.PublishResult) {
	ctx := context.Background()
	id, err := res.Get(ctx)
	if err != nil {
		panic(err)
	}
	log.Printf("Message with ID %s published", id)
}

func main() {
	publisherConfig := &ezpubsub.PublisherConfig{
		Project: "...",
		Topic: "...",
		Notifier: notify,
	}
	publisher, err := ezpubsub.NewPublisher(publisherConfig)
	if err != nil {
		panic(err)
	}

	msg := []byte("Hello world")
	publisher.Publish(msg)
}