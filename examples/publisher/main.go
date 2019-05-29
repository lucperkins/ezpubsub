package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/lucperkins/ezpubsub"
	"log"
	"time"
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
			log.Printf("Message published: (id: %s)\n", id)
		},
	}
	pub, err := ezpubsub.NewPublisher(cfg)
	must(err)

	msgs := [][]byte{[]byte("One"), []byte("Two"), []byte("Three")}
	pub.PublishBatchSync(msgs)

	userEvent := struct {
		ID        int64             `json:"id"`
		Timestamp int64             `json:"timestamp"`
		Data      map[string]string `json:"data"`
	}{
		ID:        543678,
		Timestamp: time.Now().Unix(),
		Data: map[string]string{
			"user":   "tonydanza123",
			"action": "change_username",
		},
	}

	err = pub.PublishObject(userEvent)
	must(err)

	s := fmt.Sprintf("The time now is %s", time.Now().Format("3:04PM"))

	fmt.Printf("Publishing: %s\n", s)

	pub.PublishString(s)
}
