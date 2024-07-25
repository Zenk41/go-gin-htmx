package firebase

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

func Firestore(serviceKeyFile, projectID string) (*firestore.Client, error) {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(serviceKeyFile))
	if err != nil {
		log.Fatalf("Failed to create a client: %v", err)
	}

	// be sure to close the client when done
	// defer client.Close()
	return client, err
}
