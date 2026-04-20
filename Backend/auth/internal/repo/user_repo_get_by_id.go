package repo

import (
	"auth/internal/domain/user"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

func (r *UserRepoImpl) GetByID(ctx context.Context, db *gorm.DB, userID string) (*user.UserEntity, error) {
	// 读链路：优先缓存，未命中回源 DB，并写回缓存
	cachedUser, hit, err := r.cache.Get(ctx, userID)
	if err != nil {
		r.logger.Warn("get user cache failed, fallback to db", slog.String("user_id", userID), slog.String("error", err.Error()))
	} else if hit {
		return cachedUser, nil
	}

	u, err := r.inner.GetByID(ctx, db, userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		r.cache.SetNull(ctx, userID)
		return nil, nil
	}
	if err := r.cache.Set(ctx, u); err != nil {
		r.logger.Warn("set user cache after db load failed", slog.String("user_id", userID), slog.String("error", err.Error()))
	}
	return u, nil
}

// GetByID 按主键查询用户，未命中返回 (nil, nil)
func (r *UserRepoInner) GetByID(ctx context.Context, db *gorm.DB, userID string) (*user.UserEntity, error) {
	var u user.UserEntity
	err := db.WithContext(ctx).Where("user_id = ?", userID).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user by id failed: %w", err)
	}
	return &u, nil
}
