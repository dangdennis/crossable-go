package db

var client *PrismaClient

// Client returns a singleton db client.
// It initializes the connection once if not connected yet.
func Client() *PrismaClient {
	if client != nil {
		return client
	}

	// Connect to the db via prisma
	client = NewClient()
	err := client.Connect()
	if err != nil {
		panic(err)
	}

	return client
}
