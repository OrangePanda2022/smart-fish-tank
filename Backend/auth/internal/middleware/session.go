package middleware

import (
	"auth/internal/util"
	"strings"

	"github.com/gin-gonic/gin"
)

// SessionUser 解析会话 Cookie 并把 user_id 注入上下文；会话不存在时不中断请求
func SessionUser(store util.SessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		if store == nil {
			c.Next()
			return
		}
		sid := readSessionID(c)
		if sid == "" {
			c.Next()
			return
		}
		sess, err := store.Get(c.Request.Context(), sid)
		if err != nil || sess == nil || strings.TrimSpace(sess.UserID) == "" {
			c.Next()
			return
		}
		c.Set(ContextUserIDKey, sess.UserID)
		c.Next()
	}
}

func readSessionID(c *gin.Context) string {
	sid, err := c.Cookie(SessionCookieName)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(sid)
}
