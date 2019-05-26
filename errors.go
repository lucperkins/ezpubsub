package ezpubsub

import "errors"

var (
	ErrNoProjectSpecified      = errors.New("no project specified")
	ErrNoTopicSpecified        = errors.New("no t specified")
	ErrNoSubscriptionSpecified = errors.New("no subscription specified")
	ErrNoListenerSpecified     = errors.New("no Listener function specified")
)
