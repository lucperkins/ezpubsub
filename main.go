package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"flag"
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

var process string

func init() {
	flag.StringVar(&process, "process", "subscriber", "The type of process to start")
}

func main() {
	flag.Parse()

	if process == "subscriber" {
		sub, err := messaging.NewSubscriber(project, topic, subscription, receive)
		panicOnErr(err)
		sub.Start()
	}

	if process == "publisher" {
		pub, err := messaging.NewPublisher(project, topic)
		panicOnErr(err)
		pub.Publish(context.Background(), []byte("Hello world"))
	}
}
