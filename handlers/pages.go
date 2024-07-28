package handlers

import (
	"net/http"

	"time"

	"firebase.google.com/go/auth"
	"github.com/Zenk41/go-gin-htmx/api"
	"github.com/Zenk41/go-gin-htmx/models"
	"github.com/Zenk41/go-gin-htmx/utils"
	view_auth "github.com/Zenk41/go-gin-htmx/views/auth"
	"github.com/Zenk41/go-gin-htmx/views/components"
	"github.com/Zenk41/go-gin-htmx/views/home"
	"github.com/gin-gonic/gin"
)

type PageHandler interface {
	Home(ctx *gin.Context)
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type pageHandler struct {
	userRepo     models.UserRepository
	taskRepo     models.TaskRepository
	firebaseApi  api.FirebaseApi
	firebaseAuth *auth.Client
}

func NewPageHandler(userRepo models.UserRepository,
	taskRepo models.TaskRepository,
	firebaseApi api.FirebaseApi,
	firebaseAuth *auth.Client) PageHandler {
	return &pageHandler{
		userRepo:     userRepo,
		taskRepo:     taskRepo,
		firebaseApi:  firebaseApi,
		firebaseAuth: firebaseAuth,
	}
}

func (ph *pageHandler) Home(ctx *gin.Context) {
	formattedDate := utils.GetTodayDate()

	userId, err := CookieAuth(ctx, ph.firebaseAuth)
    if err != nil {
        Render(ctx, home.Index(models.User{}, components.Alert("warning", err.Error()), formattedDate, components.Tasks([]models.Task{}, nil)))
        return
    }

	user, err := ph.userRepo.GetUser(ctx, userId)
	if err != nil {
		Render(ctx, home.Index(models.User{}, components.Alert("error", "Failed to get user"), formattedDate, components.Tasks([]models.Task{}, nil)))
		return
	}
	date, err := time.Parse("2006-01-02", formattedDate)
	if err != nil {
		// Render the page with a logged-out state if user retrieval fails
		Render(ctx, home.Index(models.User{}, components.Alert("error", err.Error()), formattedDate, components.Tasks([]models.Task{}, nil)))
		return
	}
	//Get task from repo
	tasks, err := ph.taskRepo.GetTasksByDate(ctx, user.UserID, date)
	if err != nil {
		// Render the page with a logged-out state if user retrieval fails
		Render(ctx, home.Index(models.User{}, components.Alert("error", err.Error()), formattedDate, components.Tasks([]models.Task{}, nil)))
		return
	}

	// Render the page with the logged-in state and the user data
	Render(ctx, home.Index(*user, nil, formattedDate, components.Tasks(*tasks, nil)))
}

func (ph *pageHandler) Login(ctx *gin.Context) {
	userId, errC := CookieAuth(ctx, ph.firebaseAuth)
	if errC == nil || userId != "" {
		ctx.Redirect(http.StatusFound, "/")
		return
	}
	Render(ctx, view_auth.Login(nil))
}
func (ph *pageHandler) Register(ctx *gin.Context) {
	userId, errC := CookieAuth(ctx, ph.firebaseAuth)
	if errC == nil || userId != "" {
		ctx.Redirect(http.StatusFound, "/")
		return
	}
	Render(ctx, view_auth.Register(nil))
}
