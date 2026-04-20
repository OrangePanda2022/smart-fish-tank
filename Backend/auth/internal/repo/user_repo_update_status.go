package repo

import (
	"auth/internal/domain/user"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// UpdateStatusAndTokenVersion 更新状态后删除缓存，确保状态变更立即生效
func (r *UserRepoImpl) UpdateStatus(ctx context.Context, db *gorm.DB, userID, status string) error {
	if err := r.inner.UpdateStatus(ctx, db, userID, status); err != nil {
		return err
	}
	r.cache.Delete(ctx, userID)
	return nil
}

// UpdateStatusAndTokenVersion 更新账户状态并递增 token_version
func (r *UserRepoInner) UpdateStatus(ctx context.Context, db *gorm.DB, userID, status string) error {
	// 构造更新结构体
	updateData := user.UserEntity{
		UserStatus:  status,
		UserRevoked: time.Now(),
	}

	res := db.WithContext(ctx).Model(&user.UserEntity{}).
		Where("user_id = ?", userID).
		Updates(updateData)

	if res.Error != nil {
		return fmt.Errorf("update user status failed: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
