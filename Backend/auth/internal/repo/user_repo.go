package repo

import (
	"auth/internal/cache"
	"auth/internal/domain/user"
	"context"
	"log/slog"

	"gorm.io/gorm"
)

// Repository 约束用户聚合的数据访问行为，Service 仅依赖该接口
type UserRepo interface {
	Create(ctx context.Context, db *gorm.DB, u *user.UserEntity) error
	GetByID(ctx context.Context, db *gorm.DB, userID string) (*user.UserEntity, error)
	GetByEmail(ctx context.Context, db *gorm.DB, email string) (*user.UserEntity, error)
	UpdateName(ctx context.Context, db *gorm.DB, userID, newName string) error
	UpdatePassword(ctx context.Context, db *gorm.DB, userID, passwordHash string) error
	UpdateStatus(ctx context.Context, db *gorm.DB, userID, status string) error
	UpdateRoles(ctx context.Context, db *gorm.DB, userID string, newRoles user.UserRolesList) error
	UpdatePerms(ctx context.Context, db *gorm.DB, userID string, newRoles user.UserPermsList) error
	SoftDelete(ctx context.Context, db *gorm.DB, userID string) error
}

// UserRepoImpl 在仓储层统一封装，避免业务层重复代码
type UserRepoImpl struct {
	inner  UserRepo
	cache  cache.UserByIDMapper
	logger *slog.Logger
}

// UserRepoInner 是基于 GORM 的用户仓储实现
type UserRepoInner struct {
}

// NewUserRepoCached 创建带缓存能力的仓储装饰器
func NewUserRepoImpl(inner UserRepo, cache cache.UserByIDMapper, logger *slog.Logger) *UserRepoImpl {
	return &UserRepoImpl{inner: inner, cache: cache, logger: logger}
}

// NewUserRepoInner 创建用户仓储装饰器
func NewUserRepoInner() *UserRepoInner {
	return &UserRepoInner{}
}
