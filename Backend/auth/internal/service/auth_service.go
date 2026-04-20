package service

import (
	"auth/internal/audit"
	cache "auth/internal/cache"
	"auth/internal/domain/event"
	"auth/internal/domain/status"
	"auth/internal/domain/user"
	"auth/internal/dto/response"
	errorHandler "auth/internal/error"
	"auth/internal/repo"
	"auth/internal/security"
	"auth/internal/util"
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

const (
	sessionTTL                   = time.Minute
	auditActionLogin             = "login"
	auditActionLogout            = "Logou"
	auditActionTOTPVerify        = "totp_verify"
	auditActionResetPassword     = "reset_password"
	auditActionChangePassword    = "change_password"
	loginLimiterPrefix           = "login:"
	resetRequestLimiterPrefix    = "reset_request:"
	totpVerifyLimiterPrefix      = "totp_verify:"
	singleflightGetUserKeyPrefix = "get_user:"
)

type loginLookupResult struct {
	user *user.UserEntity
}

// AuthService 封装认证与令牌生命周期相关业务
type AuthServiceImpl struct {
	database      *gorm.DB
	userRepo      repo.UserRepo
	emailMap      cache.EmailUserIDMapper
	emailBloom    cache.EmailBloomFilter
	userIDMap     cache.UserByIDMapper
	revokedMap    cache.RevokedAtUserIDMapper
	tokenCache    cache.TokenCache
	sessionStore  util.SessionStore
	audit         audit.AuditLogger
	limiter       LimiterService
	clientSecret  string
	emailMapTTL   time.Duration
	revokedMapTTL time.Duration
	resetTokenTTL time.Duration
	issuer        string
	singleflight  singleflight.Group
	logger        *slog.Logger
}

func NewAuthServiceImpl(
	database *gorm.DB,
	userRepo repo.UserRepo,
	emailMap cache.EmailUserIDMapper,
	revokedMap cache.RevokedAtUserIDMapper,
	emailBloom cache.EmailBloomFilter,
	tokenCache cache.TokenCache,
	sessionStore util.SessionStore,
	audit audit.AuditLogger,
	limiter LimiterService,
	clientSecret string,
	emailMapTTL, revokedMapTTL, resetTokenTTL time.Duration,
	issuer string,
	logger *slog.Logger,
) *AuthServiceImpl {
	// AuthService 通过端口接口拼装依赖，避免与具体基础设施强耦合
	return &AuthServiceImpl{
		database:      database,
		userRepo:      userRepo,
		emailMap:      emailMap,
		revokedMap:    revokedMap,
		emailBloom:    emailBloom,
		tokenCache:    tokenCache,
		sessionStore:  sessionStore,
		audit:         audit,
		limiter:       limiter,
		clientSecret:  clientSecret,
		emailMapTTL:   emailMapTTL,
		revokedMapTTL: revokedMapTTL,
		resetTokenTTL: resetTokenTTL,
		issuer:        issuer,
		logger:        logger,
	}
}

// Register 注册新用户并返回初始令牌与 TOTP 配置
func (s *AuthServiceImpl) Register(ctx context.Context, email, password, userName, requestID string) (*response.RegisterResult, error) {
	db := s.database
	email = util.NormalizeEmail(email)
	if hasEmpty(email, password, userName) {
		return nil, errorHandler.ErrBadRequest
	}

	// Bloom 命中才做强一致检查，可减少明显无效邮箱带来的 DB 压力
	mightExist, bloomErr := s.emailBloom.MightContain(ctx, email)
	if bloomErr != nil {
		warnDegraded(s.logger, "email bloom degraded in register", slog.String("email", email), slog.String("error", bloomErr.Error()))
		mightExist = true
	}
	if mightExist {
		existed, err := s.GetUserByEmail(ctx, email)
		if err != nil {
			return nil, errorHandler.ErrInternal.WithErr(err)
		}
		if existed != nil {
			return nil, errorHandler.ErrEmailExists
		}
	}

	passwordHash, err := security.HashPassword(password)
	if err != nil {
		return nil, errorHandler.ErrInternal.WithErr(err)
	}

	secret, _, err := util.NewTOTPSecret(s.issuer, email)
	if err != nil {
		return nil, errorHandler.ErrInternal.WithErr(err)
	}

	agg := &user.UserEntity{

		UserEmail:        email,
		UserPasswordHash: passwordHash,
		UserName:         userName,
		UserRoles:        []user.UserRole{user.UserRole{Role: user.RoleUser}},
		UserStatus:       status.StatusEnabled,
		TOTPSecret:       secret,
		UserRevoked:      time.Now(),
	}
	if err := s.userRepo.Create(ctx, db, agg); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, errorHandler.ErrEmailExists.WithErr(err)
		}
		return nil, errorHandler.ErrInternal.WithErr(err)
	}
	if err := s.emailMap.Set(ctx, email, agg.UserID, s.emailMapTTL); err != nil {
		warnDegraded(s.logger, "set email map failed", slog.String("error", err.Error()))
	}
	if err := s.emailBloom.Add(ctx, email); err != nil {
		warnDegraded(s.logger, "email bloom add failed after register", slog.String("email", email), slog.String("error", err.Error()))
	}

	return &response.RegisterResult{
		Profile: toProfile(agg),
	}, nil
}

// Login 校验账号密码并签发双令牌
func (s *AuthServiceImpl) Login(ctx context.Context, email, password, requestID string) (*response.LoginResult, error) {
	email = util.NormalizeEmail(email)
	if hasEmpty(email, password) {
		return nil, errorHandler.ErrBadRequest
	}
	if !s.limiter.Allow(loginLimiterPrefix + email) {
		return nil, errorHandler.ErrTooManyRequests
	}
	if mightExist, err := s.emailBloom.MightContain(ctx, email); err == nil && !mightExist {
		return nil, errorHandler.ErrUnauthorized
	} else if err != nil {
		warnDegraded(s.logger, "email bloom degraded in login", slog.String("email", email), slog.String("error", err.Error()))
	}

	loginUser, err := s.loadLoginUser(ctx, email)
	if err != nil {
		return nil, err
	}
	if !security.CheckPassword(loginUser.UserPasswordHash, password) {
		return nil, errorHandler.ErrPasswordMismatch
	}

	sid, _, _ := s.sessionStore.NewSession(ctx, loginUser.UserID, sessionTTL)
	s.logAudit(ctx, event.EventLogin, loginUser.UserID, email, auditActionLogin, requestID)
	return &response.LoginResult{Profile: toProfile(loginUser), SessionID: sid}, nil
}

// Logout 使会话主动失效
func (s *AuthServiceImpl) Logout(ctx context.Context, uid string, requestID string) error {
	s.logAudit(ctx, event.EventLogout, uid, "", auditActionLogout, requestID)
	return s.sessionStore.Delete(ctx, uid)
}

// RequestResetPassword 表示用户发起密码重置，服务端要求后续必须进行 TOTP 校验
func (s *AuthServiceImpl) RequestResetPassword(ctx context.Context, email string) error {
	email = util.NormalizeEmail(email)
	if email == "" {
		return errorHandler.ErrBadRequest
	}
	if !s.limiter.Allow(resetRequestLimiterPrefix + email) {
		return errorHandler.ErrTooManyRequests
	}
	// 无论用户是否存在都返回成功，防止邮箱枚举攻击

	return nil
}

// VerifyTOTPForReset 校验 TOTP 成功后签发 resetToken 并加入白名单
func (s *AuthServiceImpl) VerifyTOTPForReset(ctx context.Context, email, code, requestID string) (string, error) {
	email = util.NormalizeEmail(email)
	if hasEmpty(email, code) {
		return "", errorHandler.ErrBadRequest
	}
	if !s.limiter.Allow(totpVerifyLimiterPrefix + email) {
		return "", errorHandler.ErrTooManyRequests
	}
	if mightExist, err := s.emailBloom.MightContain(ctx, email); err == nil && !mightExist {
		return "", errorHandler.ErrUnauthorized
	} else if err != nil {
		warnDegraded(s.logger, "email bloom degraded in totp verify", slog.String("email", email), slog.String("error", err.Error()))
	}
	u, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", errorHandler.ErrUnauthorized
	}
	if !util.VerifyTOTP(u.TOTPSecret, code) {
		// TOTP 校验失败属于高价值审计事件，记录到审计链路便于安全排查
		s.logAudit(ctx, event.EventTOTPFailed, u.UserID, email, auditActionTOTPVerify, requestID)
		return "", errorHandler.ErrTOTPInvalid
	}

	resetUUID, err := uuid.NewUUID()
	if err != nil {
		return "", errorHandler.ErrInternal.WithErr(err)
	}
	resetToken := resetUUID.String()
	// resetToken 采用短 TTL 且一次性消费，降低泄露后的利用窗口
	s.tokenCache.PutResetToken(ctx, resetToken, u.UserID, s.resetTokenTTL)
	return resetToken, nil
}

// ResetPassword 使用一次性 resetToken 重置密码并使旧 Refresh Token 全部失效
func (s *AuthServiceImpl) ResetPassword(ctx context.Context, resetToken, newPassword, requestID string) error {
	db := s.database
	if hasEmpty(resetToken, newPassword) {
		return errorHandler.ErrBadRequest
	}
	userID, ok := s.tokenCache.ConsumeResetToken(ctx, resetToken)
	if !ok {
		return errorHandler.ErrResetTokenInvalid
	}
	newHash, err := security.HashPassword(newPassword)
	if err != nil {
		return errorHandler.ErrInternal.WithErr(err)
	}
	if err := s.userRepo.UpdatePassword(ctx, db, userID, newHash); err != nil {
		return mapNotFoundToAuthErr(err)
	}
	s.logAudit(ctx, event.EventPasswordReset, userID, "", auditActionResetPassword, requestID)
	return nil
}

func (s *AuthServiceImpl) GetTOTPUrl(ctx context.Context, userID string) (url string, err error) {
	db := s.database
	u, err := s.userRepo.GetByID(ctx, db, userID)
	if err != nil {
		return "", errorHandler.ErrInternal.WithErr(err)
	}
	if u == nil {
		return "", errorHandler.ErrUserNotFoundByAuth
	}
	email := u.UserEmail
	secret := u.TOTPSecret
	url = util.BuildTOTPURL(secret, s.issuer, email)
	return
}

func hasEmpty(values ...string) bool {
	for _, value := range values {
		if value == "" {
			return true
		}
	}
	return false
}

func (s *AuthServiceImpl) loadLoginUser(ctx context.Context, email string) (*user.UserEntity, error) {
	key := fmt.Sprintf("%s%s", singleflightGetUserKeyPrefix, email)
	ch := s.singleflight.DoChan(key, func() (interface{}, error) {
		u, err := s.GetUserByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
		if u == nil {
			return nil, errorHandler.ErrUnauthorized
		}
		if util.NormalizeStatus(u.UserStatus) != status.StatusEnabled {
			return nil, errorHandler.ErrUserDisabled
		}
		return &loginLookupResult{user: u}, nil
	})

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-ch:
		if res.Err != nil {
			return nil, res.Err
		}
		return res.Val.(*loginLookupResult).user, nil
	}
}

func (s *AuthServiceImpl) logAudit(ctx context.Context, eventType event.EventType, userID, userMail, action, requestID string) {
	s.audit.Log(ctx, event.AuditEvent{
		Type:      eventType,
		UserID:    userID,
		UserMail:  userMail,
		Action:    action,
		RequestID: requestID,
	})
}
