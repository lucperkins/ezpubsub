package main

import (
	"cloud.google.com/go/pubsub"
	"github.com/lucperkins/ezpubsub"
	"log"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	subscriberConfig := &ezpubsub.SubscriberConfig{
		Project:      "test",
		Topic:        "test",
		Subscription: "test",
		PushEndpoint: "http://localhost:1212",
		Listener: func(msg *pubsub.Message) {
			log.Printf("Message received: (id: %s, payload: %s)\n", msg.ID, string(msg.Data))

			msg.Ack()
		},
		ErrorHandler: func(err error) {
			log.Printf("Uh oh! An error occurred: %s", err.Error())
		},
	}

	subscriber, err := ezpubsub.NewSubscriber(subscriberConfig)
	must(err)

	subscriber.Start()

}
