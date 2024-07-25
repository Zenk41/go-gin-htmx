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
	Date        time.Time `firestore:"date"`
	CreatedAt   time.Time `firestore:"created_at"`
	UpdatedAt   time.Time `firestore:"updated_at"`
}

type TaskPayload struct {
	TaskID      string    `firestore:"task_id"`
	UserID      string    `firestore:"user_id"`
	Title       string    `firestore:"title"`
	Description string    `firestore:"description"`
	Status      string    `firestore:"status"`
	Date        time.Time `firestore:"date"`
	CreatedAt   time.Time `firestore:"created_at"`
	UpdatedAt   time.Time `firestore:"updated_at"`
}

type taskRepository struct {
	client *firestore.Client
}

type TaskRepository interface {
	GetTasksByDate(ctx context.Context, userID string, date time.Time) ([]Task, error)
	GetTodayTasks(ctx context.Context, userID string) ([]Task, error)
	CreateTask(ctx context.Context, task TaskPayload) error
	GetTaskById(ctx context.Context, taskID string) (Task, error)
	DeleteTaskById(ctx context.Context, taskID string) error
	DoneAllTaskDayByDate(ctx context.Context, userID string, date time.Time) error
	EditTaskById(ctx context.Context, task TaskPayload) error
}

func NewTaskRepository(client *firestore.Client) TaskRepository {
	return &taskRepository{
		client: client,
	}
}

// CreateTask creates a new task in Firestore
func (tr *taskRepository) CreateTask(ctx context.Context, task TaskPayload) error {
	_, err := tr.client.Collection("tasks").Doc(task.TaskID).Set(ctx, task)
	return err
}

// GetTasksByDate retrieves tasks by specific date
func (tr *taskRepository) GetTasksByDate(ctx context.Context, userID string, date time.Time) ([]Task, error) {
	var tasks []Task
	iter := tr.client.Collection("tasks").
		Where("user_id", "==", userID).
		Where("date", "==", date).
		Documents(ctx)

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
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// GetTodayTasks retrieves tasks for the current date
func (tr *taskRepository) GetTodayTasks(ctx context.Context, userID string) ([]Task, error) {
	today := time.Now().Truncate(24 * time.Hour)
	return tr.GetTasksByDate(ctx, userID, today)
}

// GetTaskById retrieves a task by its ID
func (tr *taskRepository) GetTaskById(ctx context.Context, taskID string) (Task, error) {
	doc, err := tr.client.Collection("tasks").Doc(taskID).Get(ctx)
	if err != nil {
		return Task{}, err
	}
	var task Task
	if err := doc.DataTo(&task); err != nil {
		return Task{}, err
	}
	return task, nil
}

// DeleteTaskById deletes a task by its ID
func (tr *taskRepository) DeleteTaskById(ctx context.Context, taskID string) error {
	_, err := tr.client.Collection("tasks").Doc(taskID).Delete(ctx)
	return err
}

// DoneAllTaskDayByDate marks all tasks for a specific day as done
func (tr *taskRepository) DoneAllTaskDayByDate(ctx context.Context, userID string, date time.Time) error {
	iter := tr.client.Collection("tasks").
		Where("user_id", "==", userID).
		Where("date", "==", date).
		Documents(ctx)

	batch := tr.client.Batch()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		batch.Update(doc.Ref, []firestore.Update{
			{Path: "status", Value: "done"},
		})
	}

	_, err := batch.Commit(ctx)
	return err
}

// EditTaskById edits a task by its ID
func (tr *taskRepository) EditTaskById(ctx context.Context, task TaskPayload) error {
	_, err := tr.client.Collection("tasks").Doc(task.TaskID).Set(ctx, task, firestore.MergeAll)
	return err
}
