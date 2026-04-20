package service

import (
	"auth/internal/domain/user"
	errorHandler "auth/internal/error"
	"auth/internal/util"
	"context"
	"log/slog"
	"time"

	"github.com/casbin/casbin/v3"
)

// toProfile 将领域对象转换为对外暴露的用户资料 DTO
func toProfile(u *user.UserEntity) user.UserProfile {
	roles := make([]string, len(u.UserRoles))
	for idx, val := range u.UserRoles {
		roles[idx] = val.Role
	}
	return user.UserProfile{
		UserID:     u.UserID,
		UserEmail:  u.UserEmail,
		UserName:   u.UserName,
		UserRoles:  util.NormalizeRole(roles),
		UserStatus: util.NormalizeStatus(u.UserStatus),
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

// getUserByEmail 先尝试 Email->UserID 缓存，再回源数据库，并在命中后回填缓存
func (s *AuthServiceImpl) GetUserByEmail(ctx context.Context, email string) (*user.UserEntity, error) {
	db := s.database

	if userID, ok, err := s.emailMap.Get(ctx, email, s.emailMapTTL); err == nil && ok {
		u, dbErr := s.userRepo.GetByID(ctx, db, userID)
		if u != nil && dbErr == nil {
			return u, nil
		}
		s.emailMap.Delete(ctx, email)
		return nil, errorHandler.ErrInternal.WithErr(err)
	} else if err != nil {
		warnDegraded(s.logger, "email map cache degraded", slog.String("email", email), slog.String("error", err.Error()))
	}

	u, err := s.userRepo.GetByEmail(ctx, db, email)
	if err != nil {
		return nil, errorHandler.ErrInternal.WithErr(err)
	}
	if u == nil {
		return nil, nil
	}
	if err := s.emailMap.Set(ctx, email, u.UserID, s.emailMapTTL); err != nil {
		warnDegraded(s.logger, "write email map cache failed", slog.String("error", err.Error()))
	}
	if err := s.emailBloom.Add(ctx, email); err != nil {
		warnDegraded(s.logger, "email bloom add failed after db hit", slog.String("email", email), slog.String("error", err.Error()))
	}
	return u, nil
}

func (s *AuthServiceImpl) GetRevokedAtByUserID(ctx context.Context, userID string) (time.Time, error) {
	db := s.database

	revokedAt, ok, err := s.revokedMap.Get(ctx, userID, s.revokedMapTTL)

	if err == nil && ok {
		return revokedAt, nil
	}

	if err != nil {
		warnDegraded(s.logger, "revoked map cache degraded", slog.String("user_id", userID), slog.String("error", err.Error()))
	}

	u, err := s.userRepo.GetByID(ctx, db, userID)
	if err != nil {
		return time.Time{}, errorHandler.ErrInternal.WithErr(err)
	}
	if u == nil {
		return time.Time{}, nil
	}
	if err := s.revokedMap.Set(ctx, userID, u.UserRevoked, s.revokedMapTTL); err != nil {
		warnDegraded(s.logger, "write revoked map cache failed", slog.String("error", err.Error()))
	}
	return u.UserRevoked, nil
}

func warnDegraded(logger *slog.Logger, msg string, attrs ...slog.Attr) {
	if logger == nil {
		logger = slog.Default()
	}
	args := make([]any, 0, len(attrs))
	for _, attr := range attrs {
		args = append(args, attr)
	}
	logger.Warn(msg, args...)
}

func syncRoles(enforcer *casbin.SyncedEnforcer, userID string, newRoles []string) error {
	// 获取当前 Casbin 中的旧角色
	oldRoles, err := enforcer.GetRolesForUser(userID)
	if err != nil {
		return err
	}

	// 将旧列表转为 map
	oldMap := make(map[string]bool)
	for _, r := range oldRoles {
		oldMap[r] = true
	}

	// 将新列表转为 map
	newMap := make(map[string]bool)
	for _, r := range newRoles {
		newMap[r] = true
	}

	// 计算需要删除的角色
	var toDelete [][]string
	for _, r := range oldRoles {
		if !newMap[r] {
			toDelete = append(toDelete, []string{userID, r})
		}
	}

	// 计算需要添加的角色
	var toAdd [][]string
	for _, r := range newRoles {
		if !oldMap[r] {
			toAdd = append(toAdd, []string{userID, r})
		}
	}

	// 执行批量操作
	if len(toDelete) > 0 {
		if _, err := enforcer.RemoveGroupingPolicies(toDelete); err != nil {
			return err
		}
	}
	if len(toAdd) > 0 {
		if _, err := enforcer.AddGroupingPolicies(toAdd); err != nil {
			return err
		}
	}

	return nil
}

// 将 string 转换成 user_roles
func StringToUserRoles(userID string, roleNames []string) user.UserRolesList {
	roles := make(user.UserRolesList, len(roleNames))
	for i, name := range roleNames {
		roles[i] = user.UserRole{
			UserID: userID,
			Role:   name,
		}
	}
	return roles
}

// 将 string 转换成 user_perms
func StringToUserPerms(userID string, permNames []string) user.UserPermsList {
	roles := make(user.UserPermsList, len(permNames))
	for i, name := range permNames {
		roles[i] = user.UserPerm{
			UserID: userID,
			Perm:   name,
		}
	}
	return roles
}
