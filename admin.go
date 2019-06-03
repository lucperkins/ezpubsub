package ezpubsub

type (
	// A simple administrative interface for Pub/Sub projects.
	Admin struct {
		client *client
	}
)

// Create a new Admin, specifying the Google Cloud project name.
func NewAdmin(project string) (*Admin, error) {
	client, err := newClient(project)
	if err != nil {
		return nil, err
	}

	return &Admin{
		client: client,
	}, nil
}

// List all current topics under the specified project.
func (a *Admin) ListTopics() ([]string, error) {
	return a.client.listTopics()
}

// Checks is a topic already exists.
func (a *Admin) TopicExists(topicName string) (bool, error) {
	return a.client.topicExists(topicName)
}

// Lists all current subscriptions.
func (a *Admin) ListSubscriptions() ([]string, error) {
	return a.client.listSubscriptions()
}

// Deletes a specified subscription.
func (a *Admin) DeleteSubscription(subscription string) error {
	return a.client.deleteSubscription(subscription)
}