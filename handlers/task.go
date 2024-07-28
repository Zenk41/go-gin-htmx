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
	GetTasksByDate(ctx *gin.Context)
	DeleteTaskById(ctx *gin.Context)
	DoneTaskById(ctx *gin.Context)
	DoneAllTaskDayByDate(ctx *gin.Context)
	EditTaskById(ctx *gin.Context)
	EditTaskModal(ctx *gin.Context)
	DeleteTaskModal(ctx *gin.Context)
}

type taskHandler struct {
	taskRepo     models.TaskRepository
	userRepo     models.UserRepository
	firebaseAuth *auth.Client
}

func NewTaskHandler(taskRepo models.TaskRepository, userRepo models.UserRepository, firebaseAuth *auth.Client) TaskHandler {
	return &taskHandler{
		taskRepo:     taskRepo,
		userRepo:     userRepo,
		firebaseAuth: firebaseAuth,
	}
}

// CreateNewTask handles the creation of a new task
func (th *taskHandler) CreateNewTask(ctx *gin.Context) {
	userId, errC := CookieAuth(ctx, th.firebaseAuth)
	if errC != nil || userId == "" {
		ctx.Redirect(http.StatusUnauthorized, "/login")
		return
	}
	var task models.TaskPayload

	task.Title = ctx.PostForm("title")
	task.Description = ctx.PostForm("description")
	dateStr := ctx.PostForm("date-task")
	if dateStr == "" {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/")
		return
	}
	task.Date = date

	task.TaskID = "task" + time.Now().Format("20060102150405")
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	task.UserID = userId

	if err := th.taskRepo.CreateTask(context.Background(), task); err != nil {
		Render(ctx, home.Index(models.User{}, components.Alert("error", err.Error()), dateStr, components.Tasks([]models.Task{}, nil)))
		return
	}

	user, err := th.userRepo.GetUser(context.Background(), userId)
	if err != nil {
		Render(ctx, home.Index(models.User{}, components.Alert("error", err.Error()), dateStr, components.Tasks([]models.Task{}, nil)))
		return
	}

	tasks, err := th.taskRepo.GetTasksByDate(ctx, userId, date)
	if err != nil {
		Render(ctx, home.Index(models.User{}, components.Alert("error", err.Error()), dateStr, components.Tasks([]models.Task{}, nil)))
		return
	}

	Render(ctx, home.Index(*user, components.Alert("success", "new task has been created"), dateStr, components.Tasks(*tasks, nil)))
}

// DeleteTaskById handles deleting a task by its ID
func (th *taskHandler) DeleteTaskById(ctx *gin.Context) {
	userId, errC := CookieAuth(ctx, th.firebaseAuth)
	if errC != nil || userId == "" {
		ctx.Redirect(http.StatusUnauthorized, "/login")
		return
	}
	taskID := ctx.Param("id")
	task, err := th.taskRepo.GetTaskById(ctx, taskID)
	if err != nil {
		ctx.Redirect(http.StatusFound, "/")
		return
	}

	if userId != task.UserID {
		Render(ctx, components.Task(*task, components.Alert("error", "error : You are not the owner of this task")))
		return
	}

	if err := th.taskRepo.DeleteTaskById(context.Background(), taskID); err != nil {
		th.GetTasksByDate(ctx)
		return
	}

	tasks, err := th.taskRepo.GetTasksByDate(ctx, userId, task.Date)
	if err != nil {
		Render(ctx, components.Tasks([]models.Task{}, components.Alert("error", "error : Failed to get tasks")))
		return
	}
	Render(ctx, components.Tasks(*tasks, components.Alert("success", "Task deleted successfully")))
}

func (th *taskHandler) DoneTaskById(ctx *gin.Context) {
	userId, errC := CookieAuth(ctx, th.firebaseAuth)
	if errC != nil || userId == "" {
		ctx.Redirect(http.StatusUnauthorized, "/login")
		return
	}

	task, err := th.taskRepo.GetTaskById(ctx, ctx.Param("id"))
	if err != nil {
		ctx.Redirect(http.StatusUnauthorized, "/login")
		return
	}

	if userId != task.UserID {
		th.GetTasksByDate(ctx)
		return
	}

	err = th.taskRepo.DoneTaskById(ctx, userId, task.TaskID)
	if err != nil {
		th.GetTasksByDate(ctx)
		return
	}
	date, err := time.Parse("2006-01-02",ctx.Query("date"))
	if err != nil {
		th.GetTasksByDate(ctx)
		ctx.JSON(http.StatusBadRequest,err.Error())
		return
	}

	tasks, err := th.taskRepo.GetTasksByDate(ctx, userId, date)
	if err != nil {
		th.GetTasksByDate(ctx)
		return
	}
	Render(ctx, components.Tasks(*tasks, components.Alert("success", "task by date")))
}

// DoneAllTaskDayByDate marks all tasks for a specific day as done
func (th *taskHandler) DoneAllTaskDayByDate(ctx *gin.Context) {
	userId, errC := CookieAuth(ctx, th.firebaseAuth)
	if errC != nil || userId == "" {
		ctx.Redirect(http.StatusUnauthorized, "/login")
		return
	}
	dateStr := ctx.PostForm("date-task")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		th.GetTasksByDate(ctx)
		return
	}

	firebaseCookie, err := ctx.Cookie("firebase_token")
	if err != nil || firebaseCookie == "" {
		// If there's no cookie, we still render the page with a logged-out state
		Render(ctx, components.Tasks([]models.Task{}, components.Alert("error", "error: "+err.Error())))
		return
	}

	// Verify the ID token
	token, err := th.firebaseAuth.VerifyIDToken(ctx, firebaseCookie)
	if err != nil {
		// Render the page with a logged-out state if the token verification fails
		Render(ctx, components.Tasks([]models.Task{}, components.Alert("error", "error: "+err.Error())))
		return
	}

	if err := th.taskRepo.DoneAllTaskDayByDate(context.Background(), token.UID, date); err != nil {
		Render(ctx, components.Tasks([]models.Task{}, components.Alert("error", "error: Failed to mark tasks as done")))
		return
	}

	tasks, err := th.taskRepo.GetTasksByDate(ctx, token.UID, date)
	if err != nil {
		Render(ctx, components.Tasks([]models.Task{}, components.Alert("error", "error : Failed to get tasks")))
		return
	}

	Render(ctx, components.Tasks(*tasks, components.Alert("success", "success done all task")))
}

// EditTaskById handles editing a task by its ID
func (th *taskHandler) EditTaskById(ctx *gin.Context) {
	userId, errC := CookieAuth(ctx, th.firebaseAuth)
	if errC != nil || userId == "" {
		ctx.Redirect(http.StatusUnauthorized, "/login")
		return
	}

	task, err := th.taskRepo.GetTaskById(ctx, ctx.Query("id-task"))
	if err != nil {
		ctx.Redirect(http.StatusUnauthorized, "/login")
		return
	}

	if userId != task.UserID {
		Render(ctx, components.Task(*task, components.Alert("error", "error : You are not the owner of this task")))
		return
	}

	taskPayload := models.TaskPayload{
		TaskID:      task.TaskID,
		UserID:      task.UserID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Date:        task.Date,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
	taskPayload.Title = ctx.PostForm("title")
	taskPayload.Description = ctx.PostForm("description")
	taskPayload.UpdatedAt = time.Now()

	if err := th.taskRepo.EditTaskById(ctx, taskPayload); err != nil {
		Render(ctx, components.Task(*task, components.Alert("error", "error : "+err.Error())))
		return
	}
	Render(ctx, components.Task(*task, components.Alert("success", "success : Task edited successfully")))
}

func (th *taskHandler) GetTasksByDate(ctx *gin.Context) {
	dateStr := ctx.PostForm("date-task")
	date, _ := time.Parse("2006-01-02", dateStr)
	userId, errC := CookieAuth(ctx, th.firebaseAuth)
	if errC != nil || userId == "" {
		ctx.Redirect(http.StatusUnauthorized, "/login")
		return
	}
	
	tasks, err := th.taskRepo.GetTasksByDate(ctx, userId, date)
	if err != nil {
		Render(ctx, components.Tasks([]models.Task{}, components.Alert("error", "error : Failed to get tasks")))
		return
	}
	Render(ctx, components.Tasks(*tasks, components.Alert("success", "task by date")))
}

func (th *taskHandler) EditTaskModal(ctx *gin.Context) {
	userId, errC := CookieAuth(ctx, th.firebaseAuth)
	if errC != nil || userId == "" {
		ctx.Redirect(http.StatusUnauthorized, "/login")
		return
	}
	
	idTask := ctx.Query("id-task")
	if idTask == "" {
		Render(ctx, components.ModalTaskError("error: id task not found"))
		return
	}

	task, err := th.taskRepo.GetTaskById(ctx, idTask)
	if err != nil {
		Render(ctx, components.ModalTaskError("error: "+err.Error()))
		return
	}

	Render(ctx, components.ModalEdit(*task))
}

func (th *taskHandler) DeleteTaskModal(ctx *gin.Context) {
	userId, errC := CookieAuth(ctx, th.firebaseAuth)
	if errC != nil || userId == "" {
		ctx.Redirect(http.StatusUnauthorized, "/login")
		return
	}

	
	idTask := ctx.Query("id-task")
	if idTask == "" {
		Render(ctx, components.ModalTaskError("error: id task not found"))
		return
	}

	task, err := th.taskRepo.GetTaskById(ctx, idTask)
	if err != nil {
		Render(ctx, components.ModalTaskError("error: "+err.Error()))
		return
	}

	Render(ctx, components.ModalDelete(*task))
	
}
