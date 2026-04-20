package service

import (
	"auth/internal/domain/user"
	"context"
)

// UserService 定义用户资料与管理能力
type UserService interface {
	GetProfile(ctx context.Context, userID string) (*user.UserProfile, error)
	UpdateName(ctx context.Context, userID, newName string) error
	DeleteUser(ctx context.Context, userID string) error
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword, requestID string) error
	SetUserStatus(ctx context.Context, operatorUserID, targetUserID string, enabled bool) error
	SetUserRole(ctx context.Context, operatorUserID, targetUserID string, newRoles []string) error
	SetUserPerm(ctx context.Context, operatorUserID, targetUserID string, newPerms []string) error
}
