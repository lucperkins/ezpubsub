package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/lucperkins/ezpubsub"
	"log"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	cfg := &ezpubsub.PublisherConfig{
		Project: "test",
		Topic:   "test",
		Notifier: func(res *pubsub.PublishResult) {
			id, _ := res.Get(context.Background())
			log.Printf("Message with ID %s published\n", id)
		},
	}
	pub, err := ezpubsub.NewPublisher(cfg)
	must(err)

	msgs := [][]byte{[]byte("One"), []byte("Two"), []byte("Three")}
	pub.PublishBatchSync(msgs)

}
