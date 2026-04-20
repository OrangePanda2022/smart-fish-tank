package middleware

import (
	errorHandler "auth/internal/error"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwks keyfunc.Keyfunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			AbortWithError(c, errorHandler.ErrUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// jwks.Keyfunc 根据 JWT 中的 kid 查找对应的公钥进行签名校验
		token, err := jwt.Parse(tokenString, jwks.Keyfunc)
		if err != nil || !token.Valid {
			AbortWithError(c, errorHandler.ErrUnauthorized)
			return
		}

		// 注入到 Gin Context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set(ContextUserIDKey, claims["sub"])
			c.Set(ContextUserScopesKey, claims["scp"])
		}

		c.Next()
	}
}
