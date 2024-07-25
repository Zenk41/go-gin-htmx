package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/Zenk41/go-gin-htmx/api"
	"github.com/Zenk41/go-gin-htmx/models"
	"github.com/Zenk41/go-gin-htmx/utils"
	view_auth "github.com/Zenk41/go-gin-htmx/views/auth"
	"github.com/Zenk41/go-gin-htmx/views/components"
	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type userHandler struct {
	repo        models.UserRepository
	apiKey      string
	domain      string
	firebaseApi api.FirebaseApi
}

func NewUserHandler(repo models.UserRepository, apiKey string, firebaseApi api.FirebaseApi, domain string) UserHandler {
	return &userHandler{
		repo:        repo,
		apiKey:      apiKey,
		firebaseApi: firebaseApi,
		domain:      domain,
	}
}

// Register handles the registration of a new user
func (h *userHandler) Register(ctx *gin.Context) {

	user := models.User{}

	user.Name = ctx.PostForm("name")
	user.Email = ctx.PostForm("email")
	user.Password = ctx.PostForm("password")

	registerResponse, err := h.firebaseApi.SignUpWithPassword(user.Name, user.Email, user.Password)
	if err != nil {
		errorCode, _ := utils.ExtractErrorCodeFromText(err.Error())
		Render(ctx, view_auth.Login(components.Alert("error", errorCode)))
		return
	}

	if err := user.EncryptPassword(user.Password); err != nil {
		Render(ctx, view_auth.Login(components.Alert("error", "Failed to encrypt password")))
		return
	}

	user.UserID = registerResponse["localId"].(string)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := h.repo.CreateUser(context.Background(), user); err != nil {
		errorCode, _ := utils.ExtractErrorCodeFromText(err.Error())
		Render(ctx, view_auth.Login(components.Alert("error", errorCode)))
		return
	}
	expireIn, _ := strconv.Atoi(registerResponse["expiresIn"].(string))
	ctx.SetCookie("firebase_token", registerResponse["idToken"].(string), expireIn, "/", "localhost", false, true)

	ctx.SetCookie("refresh_token", registerResponse["refreshToken"].(string), 86400, "/", "localhost", false, true)

	ctx.Redirect(http.StatusFound, "/")
}

// Login handles the user login
func (h *userHandler) Login(ctx *gin.Context) {
	user := models.User{}
	user.Email = ctx.PostForm("email")
	user.Password = ctx.PostForm("password")

	loginResponse, err := h.firebaseApi.SignInWithPassword(user.Email, user.Password)
	if err != nil {
		errorCode, _ := utils.ExtractErrorCodeFromText(err.Error())
		Render(ctx, view_auth.Login(components.Alert("error", errorCode)))
		return
	}

	expireIn, _ := strconv.Atoi(loginResponse["expiresIn"].(string))
	ctx.SetCookie("firebase_token", loginResponse["idToken"].(string), expireIn, "/", h.domain, false, false)

	ctx.SetCookie("refresh_token", loginResponse["refreshToken"].(string), 86400, "/", h.domain, false, false)

	ctx.Redirect(http.StatusFound, "/")
}

// Logout handles the user logout
func (h *userHandler) Logout(ctx *gin.Context) {
	ctx.SetCookie("firebase_token", "", 0, "/", h.domain, false, true)

	ctx.SetCookie("refresh_token", "", 0, "/", h.domain, false, true)
	Render(ctx, view_auth.Login(components.Alert("success", "logout success ")))
}
