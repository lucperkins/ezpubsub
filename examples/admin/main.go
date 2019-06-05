package main

import (
	"fmt"
	"github.com/lucperkins/ezpubsub"
)

const (
	projectName      = "test"
	topicName        = "test"
	subscriptionName = "test"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	admin, err := ezpubsub.NewAdmin(projectName)
	must(err)

	topics, err := admin.ListTopics()
	must(err)

	printList("Displaying topics:", topics, "topics")

	subscriptions, err := admin.ListSubscriptions()
	must(err)

	printList("Displaying subscriptions:", subscriptions, "subscriptions")

	topicExists, err := admin.TopicExists(topicName)
	must(err)
	fmt.Printf("Topic %s already exists: %t\n", topicName, topicExists)

	if topicExists {
		must(admin.DeleteTopic(topicName))
		fmt.Printf("Topic %s deleted\n", topicName)
	}

	subscriptionExists, err := admin.SubscriptionExists(subscriptionName)
	must(err)

	if subscriptionExists {
		must(admin.DeleteSubscription(subscriptionName))
		fmt.Printf("Subscription %s deleted\n", subscriptionName)
	}
}

func printList(msg string, list []string, listName string) {
	if len(list) != 0 {
		fmt.Println(msg)
		for _, item := range list {
			fmt.Println(item)
		}
	} else {
		fmt.Printf("The list %s is empty\n", listName)
	}

	fmt.Println()
}
