package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"flag"
	"fmt"
	"github.com/shear.io/shear/pkg/messaging"
	"log"
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

func notify(res *pubsub.PublishResult) {
	id, _ := res.Get(context.Background())
	log.Printf("Message with ID %s published successfully", id)
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
		cfg := &messaging.PublisherConfig{
			Project: project,
			Topic: topic,
			Notifier: notify,
		}
		pub, err := messaging.NewPublisher(cfg)
		panicOnErr(err)
		pub.Publish(context.Background(), []byte("Hello world"))
	}
}
