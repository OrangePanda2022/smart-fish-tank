package repo

import (
	"auth/internal/domain/user"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// SoftDelete 后写入空值缓存，降低对已删除用户的重复穿透查询
func (r *UserRepoImpl) SoftDelete(ctx context.Context, db *gorm.DB, userID string) error {
	if err := r.inner.SoftDelete(ctx, db, userID); err != nil {
		return err
	}
	r.cache.SetNull(ctx, userID)
	return nil
}

// SoftDelete 软删除用户记录（保留审计所需历史数据）
func (r *UserRepoInner) SoftDelete(ctx context.Context, db *gorm.DB, userID string) error {
	res := db.WithContext(ctx).Where("user_id = ?", userID).Delete(&user.UserEntity{})
	if res.Error != nil {
		return fmt.Errorf("soft delete user failed: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
