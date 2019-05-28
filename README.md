# ezpubsub

[![Documentation](https://godoc.org/github.com/lucperkins/ezpubsub?status.svg)](https://godoc.org/github.com/lucperkins/ezpubsub)

`ezpubsub` is a set of higher-level abstractions over the Go library for [Google Cloud Pub/Sub](https://cloud.google.com/pubsub/docs/). It's built for convenience and intended to cover the vast majority of use cases with minimal fuss. If your use case isn't covered, you're advised to use the [official library](https://godoc.org/cloud.google.com/go/pubsub).

## Why?

The [`cloud.google.com/go/pubsub`](https://godoc.org/cloud.google.com/go/pubsub) library is well done and complete but also fairly low level. `ezpubsub` makes it easier to do the ~80% of things you're most likely to do with the library, especially when using Google Cloud Pub/Sub in development and testing environments.

`ezpubsub` features a small API surface area and doesn't require you to ever deal with context; if your use case requires context, use the core `pubsub` library.

## Core concepts

The `ezpubsub` library gives you two core constructs, **publishers** and **subscribers**:

* Subscribers listen for messages on the specified topic and respond to each message using the provided listener function.
* Publishers publish messages to the specified topic either singly or in batches.

You can see examples on the [GoDoc](https://godoc.org/github.com/lucperkins/ezpubsub) page.

## Subscribers

**Subscribers** listen on the specified topic and apply the logic specified in the [`Listener`](https://godoc.org/github.com/lucperkins/ezpubsub#Listener) function to each incoming message. Here's an example:

```go
import (
        "cloud.google.com/go/pubsub"
        "github.com/lucperkins/ezpubsub"
        "log"
)

func main() {
        subscriberConfig := &ezpubsub.SubscriberConfig{
                Project: "my-project",
                Topic: "user-events",
                Subscription: "my-sub",
                Listener: func(msg *pubsub.Message) {
                        log.Printf("Event received: (id: %s, payload: %s)\n", msg.ID, string(msg.Data))
                        msg.Ack()
                },
        }
        subscriber, err := ezpubsub.NewSubscriber(subscriberConfig)
        if err != nil {
            // handle error
        }

        subscriber.Start()
}
```

## Publishers

**Publishers** publish messages on the specified topic and handle publishing results according to the logic specified in the [`Notifier`](https://godoc.org/github.com/lucperkins/ezpubsub#Notifier) function.

```go
import (
        "cloud.google.com/go/pubsub"
        "context"
        "github.com/lucperkins/ezpubsub"
        "time"
)

func main() {
        publisherConfig := &ezpubsub.PublisherConfig{
                Project: "my-project",
                Topic: "user-events",
                Notifier: func(res *pubsub.PublishResult) {
                        msgId, _ := res.Get(context.Background())
                        log.Printf("Message published: (id: %s)\n", id)
                },
        }

        publisher, err := ezpubsub.NewPublisher(publisherConfig)
        if err != nil {
                // handle error
        }

        // Publish bytes
        publisher.Publish([]byte("Hello world"))

        // Publish a string
        publisher.PublishString("Hello world")
        
        // Publish a JSON-serializable item
        event := struct {
                ID        int64
                Timestamp int64
                Message   string
        }{
                123456,
                time.Now.Uniz(),
                "Something happened",
        }
        err = publisher.PublishObject(event)
        if err != nil {
                // handle error
        }
}
```
