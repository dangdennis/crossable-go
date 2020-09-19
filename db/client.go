package db

var client *PrismaClient

// Client returns a singleton db client.
// It initializes the connection once if not connected yet.
func Client() (*PrismaClient, error) {
	if client != nil {
		return client, nil
	}

	// Connect to the db via prisma
	client = NewClient()
	err := client.Connect()
	if err != nil {
		return client, err
	}

	return client, nil
}
