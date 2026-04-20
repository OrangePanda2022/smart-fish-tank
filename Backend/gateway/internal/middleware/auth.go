package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

type AccessTokenParser interface {
	ParseAccessToken(ctx context.Context, accessToken string) (userID string, role string, err error)
}

// 解析 Bearer Token
func Auth(parser AccessTokenParser) gin.HandlerFunc {
	return func(c *gin.Context) {
		authz := c.GetHeader("Authorization")
		if authz == "" || !strings.HasPrefix(authz, "Bearer ") {
			c.AbortWithStatusJSON(401, gin.H{
				"msg": "Unauthorized",
			})
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(authz, "Bearer "))
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"msg": "Unauthorized",
			})
			return
		}
		userID, role, err := parser.ParseAccessToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"msg": "Unauthorized",
			})
			return
		}
		c.Set("user_id", userID)
		c.Set("user_role", role)
		c.Next()
	}
}
