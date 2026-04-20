package repo

import (
	"auth/internal/domain/user"
	"context"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

func (r *UserRepoImpl) Create(ctx context.Context, db *gorm.DB, u *user.UserEntity) error {
	// 先写主存储，再回填缓存；缓存失败只告警不影响主流程
	if err := r.inner.Create(ctx, db, u); err != nil {
		return err
	}
	if err := r.cache.Set(ctx, u); err != nil {
		r.logger.Warn("set user cache after create failed", slog.String("user_id", u.UserID), slog.String("error", err.Error()))
	}
	return nil
}

// Create 持久化用户聚合；唯一约束冲突由上层转换为业务错误
func (r *UserRepoInner) Create(ctx context.Context, db *gorm.DB, u *user.UserEntity) error {
	if err := db.WithContext(ctx).Create(u).Error; err != nil {
		return fmt.Errorf("create user failed: %w", err)
	}
	return nil
}
