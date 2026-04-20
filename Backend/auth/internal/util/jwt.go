package util

import (
	"auth/internal/domain/token"
	errorHandler "auth/internal/error"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenManager 抽象 JWT 签发与校验
type tokenManager interface {
	IssueAccessToken(userID string, tokenVersion int, clientID string, roles []string, scopes []string) (access string, accessExp time.Time, jti string, err error)
	IssueTokenRefresh(userID string, tokenVersion int, clientID string, roles []string, scopes []string) (refresh string, refreshExp time.Time, jti string, err error)
	ParseAccess(tok string) (*token.Claims, error)
	ParseRefresh(tok string) (*token.Claims, error)
	AccessTTL() time.Duration
	RefreshTTL() time.Duration
}

// JWTManager 负责双令牌签发与解析
type JWTManager struct {
	issuer        string
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewJWTManager(issuer, accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *JWTManager {
	return &JWTManager{
		issuer:        issuer,
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
	}
}

// AccessTTL 返回 Access Token 默认有效期
func (m *JWTManager) AccessTTL() time.Duration {
	return m.accessTTL
}

// RefreshTTL 返回 Refresh Token 默认有效期
func (m *JWTManager) RefreshTTL() time.Duration {
	return m.refreshTTL
}

// IssueAccessToken 生成 Access Token
func (m *JWTManager) IssueAccessToken(userID string, tokenVersion int, clientID string, roles []string, scopes []string) (access string, accessExp time.Time, jti string, err error) {
	now := time.Now()
	accessExp = now.Add(m.accessTTL)
	jtiUUID, err := uuid.NewV7()
	jti = jtiUUID.String()
	ac := token.Claims{
		Meta: token.TokenMeta{
			UserID:   userID,
			ClientID: clientID,
			Scopes:   scopes,
			// JTI:          jti,
			// TokenVersion: tokenVersion,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(accessExp),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, ac)
	access, err = at.SignedString(m.accessSecret)
	return
}

// IssueTokenPair 生成 Refresh 令牌
func (m *JWTManager) IssueTokenRefresh(userID string, tokenVersion int, clientID string, roles []string, scopes []string) (refresh string, refreshExp time.Time, jti string, err error) {
	now := time.Now()
	refreshExp = now.Add(m.refreshTTL)
	jtiUUID, e := uuid.NewV7()
	if e != nil {
		err = e
		return
	}
	jti = jtiUUID.String()

	rc := token.Claims{
		Meta: token.TokenMeta{
			UserID:   userID,
			ClientID: clientID,
			Scopes:   scopes,
			// JTI:          jti,
			// TokenVersion: tokenVersion,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			ID:        jti,
		},
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rc)
	refresh, err = rt.SignedString(m.refreshSecret)
	return
}

// ParseAccess 验证 Access Token
func (m *JWTManager) ParseAccess(tok string) (*token.Claims, error) {
	return m.parse(tok, m.accessSecret)
}

// ParseRefresh 验证 Refresh Token
func (m *JWTManager) ParseRefresh(tok string) (*token.Claims, error) {
	return m.parse(tok, m.refreshSecret)
}

// parse 统一完成 JWT 签名算法、过期时间和 token 类型校验
func (m *JWTManager) parse(tokenData string, secret []byte) (*token.Claims, error) {
	parsed, err := jwt.ParseWithClaims(tokenData, &token.Claims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errorHandler.ErrTokenInvalid
		}
		return secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errorHandler.ErrTokenExpired.WithErr(err)
		}
		return nil, errorHandler.ErrTokenInvalid.WithErr(err)
	}
	claims, ok := parsed.Claims.(*token.Claims)
	if !ok || !parsed.Valid {
		return nil, errorHandler.ErrTokenInvalid
	}
	return claims, nil
}
