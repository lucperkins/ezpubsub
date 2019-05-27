package main

import (
	"fmt"
	"github.com/lucperkins/ezpubsub"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	adminConfig := &ezpubsub.AdminConfig{
		Project: "test",
	}
	admin, err := ezpubsub.NewAdmin(adminConfig)
	must(err)

	topics, err := admin.ListTopics()
	must(err)

	if len(topics) != 0 {
		fmt.Println("Listing topics:")
		for _, topic := range topics {
			fmt.Println(topic)
		}
	} else {
		fmt.Println("No current topics")
	}
}
