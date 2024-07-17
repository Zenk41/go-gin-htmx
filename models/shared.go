package models

import (
	"context"
	"flag"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

func CreateClient(projectID string, firebaseApiKey string) (*firestore.Client, error) {
	ctx := context.Background()

	flag.StringVar(&projectID, "project", projectID, "The Google Cloud platform project id")
	flag.Parse()

	client, err := firestore.NewClient(ctx, projectID, option.WithAPIKey(firebaseApiKey))
	if err != nil {
		log.Fatalf("Failed to create a client: %v", err)

	}

	// be sure to close the client when done
	// defer client.Close()
	return client, err
}
