package handlers

import (
	"net/http"

	"firebase.google.com/go/auth"
	"github.com/Zenk41/go-gin-htmx/api"
	"github.com/Zenk41/go-gin-htmx/models"
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
	userRepo    models.UserRepository
	taskRepo    models.TaskRepository
	firebaseApi api.FirebaseApi
	firebaseAuth *auth.Client
}

func NewPageHandler(userRepo models.UserRepository,
	taskRepo models.TaskRepository,
	firebaseApi api.FirebaseApi,
	firebaseAuth *auth.Client) PageHandler {
	return &pageHandler{
		userRepo:    userRepo,
		taskRepo:    taskRepo,
		firebaseApi: firebaseApi,
		firebaseAuth: firebaseAuth,
	}
}

func (ph *pageHandler) Home(ctx *gin.Context) {
    firebaseCookie, err := ctx.Cookie("firebase_token")
    if err != nil || firebaseCookie == "" {
        // If there's no cookie, we still render the page with a logged-out state
        Render(ctx, home.Index([]models.Task{}, models.User{}, components.Alert("warning", "login to see your task")))
        return
    }

    // Verify the ID token
    token, err := ph.firebaseAuth.VerifyIDToken(ctx, firebaseCookie)
    if err != nil {
        // Render the page with a logged-out state if the token verification fails
        Render(ctx, home.Index([]models.Task{}, models.User{}, components.Alert("error", err.Error())))
        return
    }

    // Get user from the repository
    user, err := ph.userRepo.GetUser(ctx, token.UID)
    if err != nil || user == nil {
        // Render the page with a logged-out state if user retrieval fails
        Render(ctx, home.Index([]models.Task{}, models.User{}, components.Alert("error", err.Error())))
        return
    }

    // Render the page with the logged-in state and the user data
    Render(ctx, home.Index([]models.Task{}, *user, nil))
}

func (ph *pageHandler) Login(ctx *gin.Context) {
	firebase, _ := ctx.Cookie("firebase_token")
	if firebase != "" {
		ctx.Redirect(http.StatusFound, "/")
	}
	Render(ctx, view_auth.Login(nil))
}
func (ph *pageHandler) Register(ctx *gin.Context) {
	firebase, _ := ctx.Cookie("firebase_token")
	if firebase != "" {
		ctx.Redirect(http.StatusFound, "/")
	}
	Render(ctx, view_auth.Register(nil))
}
