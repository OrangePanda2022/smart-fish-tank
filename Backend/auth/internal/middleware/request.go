package middleware

import (
	"auth/internal/dto/result"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// ContextUserIDKey 为登录态用户 ID 的上下文键
	ContextUserIDKey = "user_id"
	// ContextClientIDKey 为 OAuth 客户端 ID 的上下文键
	ContextClientIDKey = "client_id"
	// ContextUserScopesKey 为登录态用户 scopes 的上下文键
	ContextUserScopesKey = "user_scopes"
	// ContextRequestIDKey 为请求链路 ID 的上下文键
	ContextRequestIDKey = "request_id"
	// SessionCookieName 为浏览器会话 Cookie 名称
	SessionCookieName = "sid"
)

// RequestID 为每个请求注入唯一 request_id，方便审计与排障
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := newRequestID(c.GetHeader("X-Request-Id"))
		c.Set(ContextRequestIDKey, requestID)
		c.Header("X-Request-Id", requestID)
		c.Next()
	}
}

// AbortWithError 用于中间件统一返回错误
func AbortWithError(c *gin.Context, err error) {
	c.Abort()
	result.Error(c, err)
}

func newRequestID(headerValue string) string {
	if headerValue != "" {
		return headerValue
	}
	return uuid.NewString()
}
