package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// StructuredLogger logs a gin HTTP request in JSON format using logrus.
func StructuredLogger() gin.HandlerFunc {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return func(ctx *gin.Context) {
		start := time.Now() // Start timer
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		// Process request
		ctx.Next()

		// Create log entry fields
		fields := logrus.Fields{
			"time":        time.Now().Format(time.RFC3339),
			"client_ip":   ctx.ClientIP(),
			"method":      ctx.Request.Method,
			"path":        path,
			"proto":       ctx.Request.Proto,
			"status_code": ctx.Writer.Status(),
			"latency":     time.Since(start).String(),
			"body_size":   ctx.Writer.Size(),
		}

		if raw != "" {
			fields["path"] = path + "?" + raw
		}

		if len(ctx.Errors) > 0 {
			fields["error_message"] = ctx.Errors.String()
			logger.WithFields(fields).Error("Request failed")
		} else {
			logger.WithFields(fields).Info("Request succeeded")
		}
	}
}
