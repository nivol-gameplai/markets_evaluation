package fast_markets_operations

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
)

// [END admin_import]

func CreateClient(ctx context.Context) *firestore.Client {
	// Sets your Google Cloud Platform project ID.
	projectID := "braided-triode-232313"
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

func closeClient(client *firestore.Client) {
	err := client.Close()
	if err != nil {
	}
}
