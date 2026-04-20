package repo

import (
	"auth/internal/domain/user"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// GetByEmail 仍走底层仓储，避免在装饰层重复维护多套键模型
func (r *UserRepoImpl) GetByEmail(ctx context.Context, db *gorm.DB, email string) (*user.UserEntity, error) {
	return r.inner.GetByEmail(ctx, db, email)
}

// GetByEmail 按邮箱查询用户，未命中返回 (nil, nil)
func (r *UserRepoInner) GetByEmail(ctx context.Context, db *gorm.DB, email string) (*user.UserEntity, error) {
	var u user.UserEntity
	err := db.WithContext(ctx).Where("user_email = ?", email).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by email failed: %w", err)
	}
	return &u, nil
}
