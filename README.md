# ezpubsub

[![Documentation](https://godoc.org/github.com/lucperkins/ezpubsub?status.svg)](https://godoc.org/github.com/lucperkins/ezpubsub)

`ezpubsub` is a set of higher-level abstractions over the Go library for [Google Cloud Pub/Sub](https://cloud.google.com/pubsub/docs/). It's built for convenience and intended to cover the vast majority of use cases with minimal fuss. If your use case isn't covered, you're advised to use the [official library](https://godoc.org/cloud.google.com/go/pubsub).

## Why?

The [`cloud.google.com/go/pubsub`](https://godoc.org/cloud.google.com/go/pubsub) library is well done and complete but also fairly low level. `ezpubsub` makes it easier to do the ~80% of things you're most likely to do with the library, especially when using Google Cloud Pub/Sub. It also doesn't require you to ever deal with context.

## Core concepts

The `ezpubsub` library gives you two core constructs, **publishers** and **subscribers**:

* Subscribers listen for messages on the specified topic and respond to each message using the provided listener function.
* Publishers publish messages to the specified topic either singly or in batches.

You can see examples on the [GoDoc](https://godoc.org/github.com/lucperkins/ezpubsub) page.
