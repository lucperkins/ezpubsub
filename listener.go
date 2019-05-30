package ezpubsub

import "cloud.google.com/go/pubsub"

// A subscriber listener function that does nothing but ack each message. Useful in such situations where you need
// to "wind through" outstanding messages without processing them.
func SimpleAckListener(msg *pubsub.Message) {
	msg.Ack()
}
