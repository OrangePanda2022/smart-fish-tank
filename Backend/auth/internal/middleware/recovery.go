package middleware

import (
	"auth/internal/dto/result"
	errorHandler "auth/internal/error"
	"log/slog"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery 捕获 panic 并返回统一错误，避免服务崩溃
func Recovery(logger *slog.Logger) gin.HandlerFunc {
	l := safeLogger(logger)
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				l.Error("panic recovered",
					slog.Any("panic", rec),
					slog.String("stack", string(debug.Stack())),
				)
				result.Error(c, errorHandler.ErrInternal)
				c.Abort()
			}
		}()
		c.Next()
	}
}

func safeLogger(logger *slog.Logger) *slog.Logger {
	if logger != nil {
		return logger
	}
	return slog.Default()
}
