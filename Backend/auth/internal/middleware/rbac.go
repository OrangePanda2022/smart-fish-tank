package middleware

import (
	errorHandler "auth/internal/error"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
)

func RBACMiddleware(enforcer *casbin.SyncedEnforcer, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, exists := c.Get(ContextUserIDKey)
		if exists {
			AbortWithError(c, errorHandler.ErrUnauthorized)
			return
		}

		// 角色校验
		allowed, err := enforcer.HasRoleForUser(userID.(string), requiredRole)
		if err != nil || !allowed {
			AbortWithError(c, errorHandler.ErrUnauthorized)
			return
		}

		// 注入到上下文
		c.Set("userID", userID)
		c.Next()
	}
}
