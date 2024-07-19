package handlers

import (
	utils_val "github.com/Zenk41/go-gin-htmx/utils/validate"
	views_val "github.com/Zenk41/go-gin-htmx/views/validates"
	"github.com/gin-gonic/gin"
)

type ValidateHeader interface {
	ValidateEmailHandler(ctx *gin.Context) error
	ValidatePasswordHandler(ctx *gin.Context) error
}

func ValidateEmailHandler(ctx *gin.Context) error {
	email := ctx.PostForm("email")
	valid := utils_val.IsEmailValid(email)
	return Render(ctx, views_val.ValidateEmail(valid))
}


func ValidatePasswordHandler(ctx *gin.Context) error {
	password := ctx.PostForm("password")
	valid, message := utils_val.IsPasswordValid(password)
	return Render(ctx, views_val.ValidatePassword(valid, message))
}