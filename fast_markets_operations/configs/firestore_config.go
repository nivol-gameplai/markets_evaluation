package configs

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
	"os"
)

// [END import]

// ==================================================================
// https://firebase.google.com/docs/admin/setup
// ==================================================================
// CreateClient /*
//initialization of the firestore client
//by reading project id as environment variable
//INPUTS:
//	context
//OUTPUTS:
//	firestore client
//ATTENTION: defer the client on the calling function/	*/
func CreateClient(ctx context.Context) *firestore.Client {
	// Sets your Google Cloud Platform project ID.
	//projectID := "braided-triode-232313"
	projectID := os.Getenv("GOOGLE_PROJECT_ID")
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

func CloseClient(client *firestore.Client) {
	err := client.Close()
	if err != nil {
	}
}
