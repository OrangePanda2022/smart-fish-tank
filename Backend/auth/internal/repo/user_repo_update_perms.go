package repo

import (
	"auth/internal/domain/user"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// UpdatePerms 更新角色后删除缓存，避免继续使用旧角色
func (r *UserRepoImpl) UpdatePerms(ctx context.Context, db *gorm.DB, userID string, newPerms user.UserPermsList) error {
	if err := r.inner.UpdatePerms(ctx, db, userID, newPerms); err != nil {
		return err
	}
	r.cache.Delete(ctx, userID)
	return nil
}

// UpdatePerms 更新用户角色
func (r *UserRepoInner) UpdatePerms(ctx context.Context, db *gorm.DB, userID string, newPerms user.UserPermsList) error {
	// 构造更新结构体
	updateData := user.UserEntity{
		UserPerms:   newPerms,
		UserRevoked: time.Now(),
	}
	res := db.WithContext(ctx).Model(&user.UserEntity{}).
		Where("user_id = ?", userID).
		Updates(updateData)

	if res.Error != nil {
		return fmt.Errorf("update user roles failed: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
