package main

import (
	"cloud.google.com/go/pubsub"
	"fmt"
	"github.com/lucperkins/ezpubsub"
	"strings"
)

func saveToDb(db *DB) func(*pubsub.Message) {
	return func(msg *pubsub.Message) {
		data := string(msg.Data)
		db.save(data)

		logMsg := fmt.Sprintf("Current DB contents: [%s]\n", strings.Join(db.data, ", "))
		fmt.Printf(logMsg)

		msg.Ack()
	}
}

func main() {
	db := NewDB()

	cfg := &ezpubsub.SubscriberConfig{
		Project:      "test",
		Topic:        "test",
		Subscription: "my-sub",
		Listener:     saveToDb(db),
	}
	sub, err := ezpubsub.NewSubscriber(cfg)
	must(err)

	sub.Start()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
