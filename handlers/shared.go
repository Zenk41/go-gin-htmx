package handlers

import (
	
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
