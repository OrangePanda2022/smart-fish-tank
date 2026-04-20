package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenType 表示 JWT 的用途
type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

// TokenPair 表示登录态双令牌
type TokenPair struct {
	AccessToken  AccessTokenResult  `json:"access_token"`
	RefreshToken RefreshTokenResult `json:"refresh_token"`
}

// AccessTokenResult 表示仅换发访问令牌的结果
type AccessTokenResult struct {
	AccessToken  string    `json:"access_token"`
	AccessExpire time.Time `json:"access_expire"`
}

// RefreshTokenResult 表示仅换发访问令牌的结果
type RefreshTokenResult struct {
	RefreshToken  string    `json:"refresh_token"`
	RefreshExpire time.Time `json:"refresh_expire"`
}

// Claims 统一承载访问令牌与刷新令牌信息
type Claims struct {
	Meta TokenMeta `json:"meta"`
	jwt.RegisteredClaims
}

type TokenMeta struct {
	ClientID string   `json:"client_id"`
	UserID   string   `json:"uid"`
	Scopes   []string `json:"scopes"`
	// JTI          string
	// TokenVersion int
}
