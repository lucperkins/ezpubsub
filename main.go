package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/shear.io/shear/pkg/messaging"
)

const (
	project = "shear-dev"
	topic = "test-topic"
	subscription = "my-sub"
)

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func receive(_ context.Context, msg *pubsub.Message) {
	fmt.Printf("Message received: %s", string(msg.Data))
	msg.Ack()
}

func main() {
	sub, err := messaging.NewSubscriber(project, topic, subscription, receive)
	panicOnErr(err)

	go func() {
		sub.Start()
	}()

	pub, err := messaging.NewPublisher(project, topic)
	panicOnErr(err)

	pub.Publish(context.Background(), []byte("Hello world"))
}
