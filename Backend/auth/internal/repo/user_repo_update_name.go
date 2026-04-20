package repo

import (
	"auth/internal/domain/user"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// UpdateName 更新后主动删缓存，采用“删除而非更新”避免并发脏写
func (r *UserRepoImpl) UpdateName(ctx context.Context, db *gorm.DB, userID, newName string) error {
	if err := r.inner.UpdateName(ctx, db, userID, newName); err != nil {
		return err
	}
	r.cache.Delete(ctx, userID)
	return nil
}

// UpdateName 修改昵称，RowsAffected=0 视为用户不存在
func (r *UserRepoInner) UpdateName(ctx context.Context, db *gorm.DB, userID, newName string) error {
	res := db.WithContext(ctx).Model(&user.UserEntity{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{"user_name": newName})
	if res.Error != nil {
		return fmt.Errorf("update user name failed: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
