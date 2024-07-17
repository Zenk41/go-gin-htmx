package models

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
)

// User represents a user in the application
type User struct {
	UserID    string    `firestore:"user_id"`
	Email     string    `firestore:"email"`
	Password  string    `firestore:"password"` // Ensure this is hashed before storing
	Name      string    `firestore:"name"`
	CreatedAt time.Time `firestore:"created_at"`
	UpdatedAt time.Time `firestore:"updated_at"`
}

// CreateUser creates a new user in Firestore
func CreateUser(ctx context.Context, client *firestore.Client, user User) error {
	_, err := client.Collection("users").Doc(user.UserID).Set(ctx, user)
	return err
}

// GetUser retrieves a user by ID from Firestore
func GetUser(ctx context.Context, client *firestore.Client, userID string) (*User, error) {
	doc, err := client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		return nil, err
	}
	var user User
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
