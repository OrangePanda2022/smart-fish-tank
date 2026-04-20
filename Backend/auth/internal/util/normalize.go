package util

import (
	"auth/internal/domain/status"
	"strings"
	"time"
)

// NormalizeEmail 标准化邮箱格式（去首尾空格并转小写）
func NormalizeEmail(email string) string {
	return normalizeString(email)
}

// NormalizeRole 标准化角色集合（原地修改）
func NormalizeRole(roles []string) []string {
	return normalizeStringSlice(roles)
}

// NormalizeScopes 标准化权限集合（原地修改）
func NormalizeScopes(scopes []string) []string {
	return normalizeStringSlice(scopes)
}

// NormalizeStatus 标准化用户状态；空值回落为启用态
func NormalizeStatus(s string) string {
	normalized := normalizeString(s)
	if normalized == "" {
		return status.StatusEnabled
	}
	return normalized
}

// NormalizeTTL 标准化 TTL，避免出现小于 1 秒的非正值
func NormalizeTTL(ttl time.Duration) time.Duration {
	if ttl == 0 {
		return 0
	}
	if ttl < 0 {
		return time.Second
	}
	if ttl < time.Second {
		return time.Second
	}
	return ttl
}

func normalizeStringSlice(values []string) []string {
	for i := range values {
		values[i] = normalizeString(values[i])
	}
	return values
}

func normalizeString(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
