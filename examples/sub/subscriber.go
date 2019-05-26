package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/lucperkins/ezpubsub"
	"log"
)

func processMessage(_ context.Context, msg *pubsub.Message) {
	log.Printf("Message received with an ID of %s and the following payload: %s", msg.ID, string(msg.Data))
	msg.Ack()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	subscriberConfig := &ezpubsub.SubscriberConfig{
		Project:      "...",
		Topic:        "...",
		Subscription: "...",
		Listener:     processMessage,
	}
	subscriber, err := ezpubsub.NewSubscriber(subscriberConfig)
	if err != nil {
		panic(err)
	}

	subscriber.Start()
}
