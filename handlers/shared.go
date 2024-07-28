package handlers

import (
	"errors"

	"firebase.google.com/go/auth"
	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HTTPHandler is a type for handlers that return an error
type HTTPHandler func(c *gin.Context) error

// Make wraps an HTTPHandler into a gin.HandlerFunc
func Make(h HTTPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h(c); err != nil {
			logrus.Error("HTTP handle error", "err", err, "path", c.Request.URL.Path)
		}
	}
}

// Render renders a templ.Component using the gin.Context
func Render(c *gin.Context, t templ.Component) error {
	return t.Render(c.Request.Context(), c.Writer)
}

// CookieAuth retrieves the UID from the cookie and verifies it with Firebase Auth
func CookieAuth(ctx *gin.Context, auth *auth.Client) (string, error) {

	firebaseCookie, err := ctx.Cookie("firebase_token")
	if err != nil || firebaseCookie == "" {
		if err != nil {
			return "", errors.New("login to see your task: " + err.Error())
		}

		return "", errors.New("login to see your task: cookie is empty")
	}

	token, err := auth.VerifyIDToken(ctx, firebaseCookie)
	if err != nil {

		return "", errors.New("token verification failed: " + err.Error())
	}

	return token.UID, nil
}
