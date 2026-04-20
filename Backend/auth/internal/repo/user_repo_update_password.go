package repo

import (
	"auth/internal/domain/user"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// UpdatePasswordAndTokenVersion 更新凭据后立即删缓存，确保认证相关字段强一致
func (r *UserRepoImpl) UpdatePassword(ctx context.Context, db *gorm.DB, userID, passwordHash string) error {
	if err := r.inner.UpdatePassword(ctx, db, userID, passwordHash); err != nil {
		return err
	}
	r.cache.Delete(ctx, userID)
	return nil
}

// UpdatePasswordAndTokenVersion 原子更新密码并递增 token_version
func (r *UserRepoInner) UpdatePassword(ctx context.Context, db *gorm.DB, userID, passwordHash string) error {
	res := db.WithContext(ctx).Model(&user.UserEntity{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"user_password_hash": passwordHash,
			"user_revoked":       time.Now(),
		})
	if res.Error != nil {
		return fmt.Errorf("update password failed: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
