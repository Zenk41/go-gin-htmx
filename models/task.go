package models

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// Task represents a task in the application
type Task struct {
	TaskID      string    `firestore:"task_id"`
	UserID      string    `firestore:"user_id"`
	Title       string    `firestore:"title"`
	Description string    `firestore:"description"`
	Status      string    `firestore:"status"`
	CreatedAt   time.Time `firestore:"created_at"`
	UpdatedAt   time.Time `firestore:"updated_at"`
}

// CreateTask creates a new task in Firestore
func CreateTask(ctx context.Context, client *firestore.Client, task Task) error {
	_, err := client.Collection("tasks").Doc(task.TaskID).Set(ctx, task)
	return err
}

// GetTasks retrieves tasks for a user from Firestore
func GetTasks(ctx context.Context, client *firestore.Client, userID string) ([]*Task, error) {
	var tasks []*Task
	iter := client.Collection("tasks").Where("user_id", "==", userID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var task Task
		if err := doc.DataTo(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}
