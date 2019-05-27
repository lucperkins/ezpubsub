package ezpubsub

type (
	// A simple administrative interface for Pub/Sub projects.
	Admin struct {
		client *client
	}
)

// List all current topics under the specified project.
func (a *Admin) ListTopics() ([]string, error) {
	return a.client.listTopics()
}

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
