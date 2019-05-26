package examples

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/lucperkins/ezpubsub"
	"log"
)

func notify(res *pubsub.PublishResult) {
	ctx := context.Background()
	id, err := res.Get(ctx)
	must(err)
	log.Printf("Message with ID %s published", id)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	publisherConfig := &ezpubsub.PublisherConfig{
		Project: "...",
		Topic: "...",
		Notifier: notify,
	}
	publisher, err := ezpubsub.NewPublisher(publisherConfig)
	must(err)

	msg := []byte("Hello world")
	publisher.Publish(msg)
}