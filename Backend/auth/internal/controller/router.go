package controller

import (
	"auth/internal/middleware"
	"auth/internal/util"
	"log/slog"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

func NewRouter(handler *Handler, jwks keyfunc.Keyfunc, sessionStore util.SessionStore, redis redis.UniversalClient, httpRateLimitRPS, httpRateLimitBurst int, logger *slog.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.RequestID())
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.HTTPRateLimitLocal(rate.Limit(httpRateLimitRPS), httpRateLimitBurst, logger))
	r.Use(middleware.HTTPRateLimitRedis(redis, httpRateLimitBurst, logger))
	r.Use(middleware.SessionUser(sessionStore))

	// 所有业务路由统一挂到 /api/v1，便于版本化演进
	v1 := r.Group("/api/v1")
	{
		// Hydra login/consent 回调端点（登录同意页服务端）
		oauth2Group := v1.Group("/oauth2")
		{
			oauth2Group.GET("/login", handler.OAuth2LoginRequest)
			oauth2Group.GET("/consent", handler.OAuth2ConsentRequest)
			oauth2Group.POST("/consent", handler.OAuth2ConsentSubmit)
		}

		// 公开认证接口（无需登录态）
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", handler.Register)
			authGroup.POST("/login", handler.Login)
			authGroup.POST("/totp/verify", handler.VerifyTOTP)
			authGroup.POST("/password/reset/request", handler.RequestResetPassword)
			authGroup.POST("/password/reset/confirm", handler.ResetPassword)
		}

		// 受保护接口：先通过会话鉴权，再进入业务处理
		userProtected := v1.Group("/users")
		userProtected.Use(middleware.AuthMiddleware(jwks))
		{
			userProtected.POST("/logout", handler.Logout)
			userProtected.GET("/me", handler.GetProfile)
			userProtected.PATCH("/name", handler.UpdateName)
			userProtected.PATCH("/password", handler.ChangePassword)
			userProtected.DELETE("/me", handler.DeleteUser)
		}

		// 管理员接口：当前仅做登录态保护，具体角色校验由服务层兜底
		adminProtected := v1.Group("/admin")
		adminProtected.Use(middleware.AuthMiddleware(jwks))
		{
			adminProtected.PATCH("/users/:user_id/status", handler.AdminSetUserStatus)
			adminProtected.PATCH("/users/:user_id/role", handler.AdminSetUserRole)
		}
	}

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	return r
}
