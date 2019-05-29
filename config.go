package ezpubsub

type (
	// Publisher configuration. All fields except Notifier are mandatory.
	PublisherConfig struct {
		Project  string
		Topic    string
		Notifier Notifier
	}

	// Subscriber configuration. A Project, Topic, and Subscription are mandatory; errors are thrown if these are not
	// provided. A Listener function is optional; if none is provided, a defaultListener is used that for each message
	// received logs a simple string and acks the message. An ErrorHandler function is also optional; if none is
	// provided, errors are logged to stderr.
	SubscriberConfig struct {
		Project      string
		Topic        string
		Subscription string
		Listener     Listener
		ErrorHandler ErrorHandler
	}
)

// Validate the PublisherConfig
func (c *PublisherConfig) validate() error {
	if c.Project == "" {
		return ErrNoProjectSpecified
	}
	if c.Topic == "" {
		return ErrNoTopicSpecified
	}

	return nil
}

func (c *SubscriberConfig) validate() error {
	if c.Project == "" {
		return ErrNoProjectSpecified
	}
	if c.Topic == "" {
		return ErrNoTopicSpecified
	}
	if c.Subscription == "" {
		return ErrNoSubscriptionSpecified
	}
	if c.Listener == nil {
		c.Listener = defaultListener
	}
	if c.ErrorHandler == nil {
		c.ErrorHandler = defaultErrorHandler
	}

	return nil
}
