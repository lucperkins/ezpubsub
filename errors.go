package ezpubsub

import "errors"

var (
	ErrNoProjectSpecified      = errors.New("no project specified")
	ErrNoTopicSpecified        = errors.New("no topic specified")
	ErrNoSubscriptionSpecified = errors.New("no subscription specified")
)
