package handlers

import (
	"context"
	"net/http"
	"time"

	"firebase.google.com/go/auth"
	"github.com/Zenk41/go-gin-htmx/models"
	"github.com/Zenk41/go-gin-htmx/views/components"
	"github.com/Zenk41/go-gin-htmx/views/home"
	"github.com/gin-gonic/gin"
)

type TaskHandler interface {
	CreateNewTask(ctx *gin.Context)
	GetTodayTask(ctx *gin.Context)
	GetTaskById(ctx *gin.Context)
	GetTasksByDate(ctx *gin.Context)
	DeleteTaskById(ctx *gin.Context)
	DoneAllTaskDayByDate(ctx *gin.Context)
	EditTaskById(ctx *gin.Context)
}

type taskHandler struct {
	repo models.TaskRepository
	firebaseAuth *auth.Client
}

func NewTaskHandler(repo models.TaskRepository, firebaseAuth *auth.Client) TaskHandler {
	return &taskHandler{
		repo: repo,
		firebaseAuth: firebaseAuth,
	}
}

// CreateNewTask handles the creation of a new task
func (th *taskHandler) CreateNewTask(ctx *gin.Context) {
	
	var task models.TaskPayload

	task.Title = ctx.PostForm("title")
	task.Description = ctx.PostForm("description")
	dateStr := ctx.PostForm("datetask")

	date, _ := time.Parse("2006-01-02", dateStr)
	task.Date = date

	task.TaskID = "task" + time.Now().Format("20060102150405")
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	firebaseCookie, err := ctx.Cookie("firebase_token")
	if err != nil || firebaseCookie == "" {
		// If there's no cookie, we still render the page with a logged-out state
		Render(ctx, home.Index( models.User{}, components.Alert("warning", "login to see your task"), dateStr,components.Tasks([]models.Task{})))
		return
	}

	// Verify the ID token
	token, err := th.firebaseAuth.VerifyIDToken(ctx, firebaseCookie)
	if err != nil {
		// Render the page with a logged-out state if the token verification fails
		Render(ctx, home.Index(models.User{}, components.Alert("error", err.Error()), dateStr, components.Tasks([]models.Task{})))
		return
	}
	task.UserID = token.UID

	if err := th.repo.CreateTask(context.Background(), task); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task created successfully"+task.Date.String() + dateStr})
}

// GetTodayTask handles retrieving today's tasks for a specific user
func (th *taskHandler) GetTodayTask(ctx *gin.Context) {
	userID := ctx.Param("userID")
	tasks, err := th.repo.GetTodayTasks(context.Background(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// GetTaskById handles retrieving a task by its ID
func (th *taskHandler) GetTaskById(ctx *gin.Context) {
	
	taskID := ctx.Param("taskID")
	task, err := th.repo.GetTaskById(context.Background(), taskID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"task": task})
}

// DeleteTaskById handles deleting a task by its ID
func (th *taskHandler) DeleteTaskById(ctx *gin.Context) {
	taskID := ctx.Param("taskID")
	if err := th.repo.DeleteTaskById(context.Background(), taskID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// DoneAllTaskDayByDate marks all tasks for a specific day as done
func (th *taskHandler) DoneAllTaskDayByDate(ctx *gin.Context) {
	userID := ctx.Param("userID")
	dateStr := ctx.Param("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	if err := th.repo.DoneAllTaskDayByDate(context.Background(), userID, date); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark tasks as done"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "All tasks for the day marked as done"})
}

// EditTaskById handles editing a task by its ID
func (th *taskHandler) EditTaskById(ctx *gin.Context) {
	var task models.TaskPayload
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	taskID := ctx.Param("taskID")
	task.TaskID = taskID
	task.UpdatedAt = time.Now()

	if err := th.repo.EditTaskById(context.Background(), task); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit task"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task edited successfully"})
}


func (th *taskHandler) GetTasksByDate(ctx *gin.Context) {
	dateStr := ctx.PostForm("date-task")
	date, _ := time.Parse("2006-01-02", dateStr)
    firebaseCookie, err := ctx.Cookie("firebase_token")
    if err != nil || firebaseCookie == "" {
        // If there's no cookie, we still render the page with a logged-out state
        Render(ctx, components.Tasks([]models.Task{}))
        return
    }

    // Verify the ID token
    token, err := th.firebaseAuth.VerifyIDToken(ctx, firebaseCookie)
    if err != nil {
        // Render the page with a logged-out state if the token verification fails
        Render(ctx, components.Tasks([]models.Task{}))
        return
    }
	tasks,err := th.repo.GetTasksByDate(ctx, token.UID, date)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks"})
		return
	}
	Render(ctx, components.Tasks(*tasks))
}