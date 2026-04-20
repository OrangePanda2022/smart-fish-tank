package middleware

import (
	errorHandler "auth/internal/error"
	"strings"

	"github.com/gin-gonic/gin"
)

// RequireScopes 要求当前登录用户 scopes 中至少包含一个目标 scope，否则返回 403
func RequireScopes(scopes ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(scopes))
	for _, scope := range scopes {
		normalized := normalizeScope(scope)
		if normalized != "" {
			allowed[normalized] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		if len(allowed) == 0 {
			c.Next()
			return
		}

		v, ok := c.Get(ContextUserScopesKey)
		if !ok {
			AbortWithError(c, errorHandler.ErrUnauthorized)
			return
		}
		userScopes, ok := toStringSlice(v)
		if !ok || len(userScopes) == 0 {
			AbortWithError(c, errorHandler.ErrUnauthorized)
			return
		}
		for _, scope := range userScopes {
			normalized := normalizeScope(scope)
			if _, exists := allowed[normalized]; exists {
				c.Next()
				return
			}
		}
		AbortWithError(c, errorHandler.ErrForbidden)
	}
}

// toStringSlice 将上下文里的 scope 值转换为字符串切片
func toStringSlice(v any) ([]string, bool) {
	switch vv := v.(type) {
	case []string:
		return vv, true
	case []any:
		out := make([]string, 0, len(vv))
		for _, item := range vv {
			s, ok := item.(string)
			if !ok {
				return nil, false
			}
			out = append(out, s)
		}
		return out, true
	default:
		return nil, false
	}
}

func normalizeScope(scope string) string {
	return strings.TrimSpace(strings.ToLower(scope))
}
