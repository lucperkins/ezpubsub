package ezpubsub

import "errors"

var (
	ErrNoProjectSpecified      = errors.New("no project specified")
	ErrNoTopicSpecified        = errors.New("no topic specified")
	ErrNoSubscriptionSpecified = errors.New("no subscription specified")
)

// A function that determines how errors are handled.
type ErrorHandler = func(error)
