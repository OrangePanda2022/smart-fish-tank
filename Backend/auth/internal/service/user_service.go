package service

import (
	cache "auth/internal/cache"
	userStatus "auth/internal/domain/status"
	"auth/internal/domain/user"
	errorHandler "auth/internal/error"
	"auth/internal/repo"
	"auth/internal/security"
	"context"
	"errors"

	"github.com/casbin/casbin/v3"
	"gorm.io/gorm"
)

// UserService 负责用户资料相关业务
type UserServiceImpl struct {
	database *gorm.DB
	userRepo repo.UserRepo
	emailMap cache.EmailUserIDMapper
	enforcer *casbin.SyncedEnforcer
}

func NewUserServiceImpl(database *gorm.DB, repo repo.UserRepo, emailMap cache.EmailUserIDMapper, enforcer *casbin.SyncedEnforcer) *UserServiceImpl {
	return &UserServiceImpl{database: database, userRepo: repo, emailMap: emailMap, enforcer: enforcer}
}

// GetProfile 查询当前用户资料
func (s *UserServiceImpl) GetProfile(ctx context.Context, userID string) (*user.UserProfile, error) {
	userEntity, err := s.getExistingUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	profile := toProfile(userEntity)
	return &profile, nil
}

// UpdateName 修改昵称
func (s *UserServiceImpl) UpdateName(ctx context.Context, userID, newName string) error {
	db := s.database
	if userID == "" || newName == "" {
		return errorHandler.ErrBadRequest
	}
	if err := s.userRepo.UpdateName(ctx, db, userID, newName); err != nil {
		return mapNotFoundToAuthErr(err)
	}
	return nil
}

// DeleteUser 软删除用户并清理邮箱映射缓存
func (s *UserServiceImpl) DeleteUser(ctx context.Context, userID string) error {
	db := s.database
	u, err := s.getExistingUser(ctx, userID)
	if err != nil {
		return err
	}
	if err := s.userRepo.SoftDelete(ctx, db, userID); err != nil {
		return mapNotFoundToAuthErr(err)
	}
	s.emailMap.Delete(ctx, u.UserEmail)
	return nil
}

// ChangePassword 登录用户修改密码
func (s *UserServiceImpl) ChangePassword(ctx context.Context, userID, oldPassword, newPassword, requestID string) error {
	db := s.database
	if hasEmpty(userID, oldPassword, newPassword) {
		return errorHandler.ErrBadRequest
	}
	u, err := s.userRepo.GetByID(ctx, db, userID)
	if err != nil {
		return errorHandler.ErrInternal.WithErr(err)
	}
	if u == nil {
		return errorHandler.ErrUserNotFoundByAuth
	}
	if !security.CheckPassword(u.UserPasswordHash, oldPassword) {
		return errorHandler.ErrPasswordMismatch
	}
	hash, err := security.HashPassword(newPassword)
	if err != nil {
		return errorHandler.ErrInternal.WithErr(err)
	}
	if err := s.userRepo.UpdatePassword(ctx, db, userID, hash); err != nil {
		return mapNotFoundToAuthErr(err)
	}
	return nil
}

// SetUserStatus 管理员启用或禁用用户账号
func (s *UserServiceImpl) SetUserStatus(ctx context.Context, operatorUserID, targetUserID string, enabled bool) error {
	db := s.database
	if operatorUserID == "" || targetUserID == "" {
		return errorHandler.ErrBadRequest
	}
	if operatorUserID == targetUserID && !enabled {
		return errorHandler.ErrForbidden
	}

	_, err := s.getOperator(ctx, operatorUserID)
	if err != nil {
		return err
	}
	allowed, err := s.enforcer.Enforce(operatorUserID, "user", "setStatus")
	if err != nil || !allowed {
		return errors.New("无权执行此操作")
	}

	status := userStatus.StatusDisabled
	if enabled {
		status = userStatus.StatusEnabled
	}
	if err := s.userRepo.UpdateStatus(ctx, db, targetUserID, status); err != nil {
		return mapNotFoundToAuthErr(err)
	}

	return nil
}

// SetUserRole 管理员授予或取消目标用户管理员身份，并使目标旧令牌失效
func (s *UserServiceImpl) SetUserRole(ctx context.Context, operatorUserID, targetUserID string, newRoles []string) error {
	db := s.database
	if operatorUserID == "" || targetUserID == "" {
		return errorHandler.ErrBadRequest
	}

	// TODO 校验 role 是否合法

	_, err := s.getOperator(ctx, operatorUserID)
	if err != nil {
		return err
	}
	allowed, err := s.enforcer.Enforce(operatorUserID, "user", "setRoles")
	if err != nil || !allowed {
		return errors.New("无权执行此操作")
	}

	if err := s.userRepo.UpdateRoles(ctx, db, targetUserID, StringToUserRoles(targetUserID, newRoles)); err != nil {
		return mapNotFoundToAuthErr(err)
	}

	// casbin 多实例异步修改 roles
	if err := syncRoles(s.enforcer, targetUserID, newRoles); err != nil {
		s.enforcer.LoadPolicy()
		return err
	}

	return nil
}

// SetUserPerm 管理员授予或取消目标用户管理员权限，并使目标旧令牌失效
func (s *UserServiceImpl) SetUserPerm(ctx context.Context, operatorUserID, targetUserID string, newPerms []string) error {
	db := s.database
	if operatorUserID == "" || targetUserID == "" {
		return errorHandler.ErrBadRequest
	}

	// TODO 校验 perm 是否合法
	_, err := s.getOperator(ctx, operatorUserID)

	if err != nil {
		return err
	}

	allowed, err := s.enforcer.Enforce(operatorUserID, "user", "setPerms")
	if err != nil || !allowed {
		return errors.New("无权执行此操作")
	}

	if err := s.userRepo.UpdatePerms(ctx, db, targetUserID, StringToUserPerms(targetUserID, newPerms)); err != nil {
		return mapNotFoundToAuthErr(err)
	}

	// casbin 多实例异步修改 perms
	if err := syncRoles(s.enforcer, targetUserID, newPerms); err != nil {
		s.enforcer.LoadPolicy()
		return err
	}

	return nil
}

func (s *UserServiceImpl) getExistingUser(ctx context.Context, userID string) (*user.UserEntity, error) {
	u, err := s.userRepo.GetByID(ctx, s.database, userID)
	if err != nil {
		return nil, errorHandler.ErrInternal.WithErr(err)
	}
	if u == nil {
		return nil, errorHandler.ErrUserNotFoundByAuth
	}
	return u, nil
}

func (s *UserServiceImpl) getOperator(ctx context.Context, operatorUserID string) (*user.UserEntity, error) {
	operator, err := s.userRepo.GetByID(ctx, s.database, operatorUserID)
	if err != nil {
		return nil, errorHandler.ErrInternal.WithErr(err)
	}
	if operator == nil {
		return nil, errorHandler.ErrUnauthorized
	}
	return operator, nil
}

func mapNotFoundToAuthErr(err error) error {
	if err == gorm.ErrRecordNotFound {
		return errorHandler.ErrUserNotFoundByAuth
	}
	return errorHandler.ErrInternal.WithErr(err)
}
