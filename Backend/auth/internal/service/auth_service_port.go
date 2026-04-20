package service

import (
	"auth/internal/dto/response"
	"context"
	"time"
)

// AuthService 定义认证与账号安全相关用例
type AuthService interface {
	Register(ctx context.Context, email, password, userName, requestID string) (*response.RegisterResult, error)
	Login(ctx context.Context, email, password, requestID string) (*response.LoginResult, error)
	Logout(ctx context.Context, userID string, requestID string) error
	RequestResetPassword(ctx context.Context, email string) error
	VerifyTOTPForReset(ctx context.Context, email, code, requestID string) (string, error)
	GetTOTPUrl(ctx context.Context, userID string) (url string, err error)
	ResetPassword(ctx context.Context, resetToken, newPassword, requestID string) error
	GetRevokedAtByUserID(ctx context.Context, userID string) (time.Time, error)
}
