package repo

import (
	"auth/internal/domain/user"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// UpdateRoles 更新角色后删除缓存，避免继续使用旧角色
func (r *UserRepoImpl) UpdateRoles(ctx context.Context, db *gorm.DB, userID string, newRoles user.UserRolesList) error {
	if err := r.inner.UpdateRoles(ctx, db, userID, newRoles); err != nil {
		return err
	}
	r.cache.Delete(ctx, userID)
	return nil
}

// UpdateRoles 更新用户角色
func (r *UserRepoInner) UpdateRoles(ctx context.Context, db *gorm.DB, userID string, newRoles user.UserRolesList) error {
	// 构造更新结构体
	updateData := user.UserEntity{
		UserRoles:   newRoles,
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
