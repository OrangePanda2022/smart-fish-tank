package errorHandler

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError 是统一业务错误模型，可映射为 HTTP 状态码并携带内部错误原因
type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
	Err        error
}

func (e *AppError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// WithErr 将底层错误挂载到业务错误，方便日志排障
func (e *AppError) WithErr(err error) *AppError {
	clone := *e
	clone.Err = err
	return &clone
}

// Is 用于 errors.Is 对比同类业务错误
func (e *AppError) Is(target error) bool {
	var t *AppError
	if !errors.As(target, &t) {
		return false
	}
	return e.Code == t.Code
}

var (
	ErrBadRequest         = &AppError{Code: "BAD_REQUEST", Message: "请求参数错误", HTTPStatus: http.StatusBadRequest}
	ErrUnauthorized       = &AppError{Code: "UNAUTHORIZED", Message: "未授权访问", HTTPStatus: http.StatusUnauthorized}
	ErrForbidden          = &AppError{Code: "FORBIDDEN", Message: "无权限执行该操作", HTTPStatus: http.StatusForbidden}
	ErrNotFound           = &AppError{Code: "NOT_FOUND", Message: "资源不存在", HTTPStatus: http.StatusNotFound}
	ErrConflict           = &AppError{Code: "CONFLICT", Message: "资源冲突", HTTPStatus: http.StatusConflict}
	ErrTooManyRequests    = &AppError{Code: "TOO_MANY_REQUESTS", Message: "请求过于频繁", HTTPStatus: http.StatusTooManyRequests}
	ErrTokenInvalid       = &AppError{Code: "TOKEN_INVALID", Message: "令牌无效", HTTPStatus: http.StatusUnauthorized}
	ErrTokenExpired       = &AppError{Code: "TOKEN_EXPIRED", Message: "令牌已过期", HTTPStatus: http.StatusUnauthorized}
	ErrTokenBlacklisted   = &AppError{Code: "TOKEN_BLACKLISTED", Message: "令牌已失效", HTTPStatus: http.StatusUnauthorized}
	ErrPasswordMismatch   = &AppError{Code: "PASSWORD_MISMATCH", Message: "密码不正确", HTTPStatus: http.StatusUnauthorized}
	ErrEmailExists        = &AppError{Code: "EMAIL_EXISTS", Message: "邮箱已被注册", HTTPStatus: http.StatusConflict}
	ErrUserDeleted        = &AppError{Code: "USER_DELETED", Message: "用户已注销", HTTPStatus: http.StatusGone}
	ErrUserDisabled       = &AppError{Code: "USER_DISABLED", Message: "用户已被禁用", HTTPStatus: http.StatusForbidden}
	ErrTOTPInvalid        = &AppError{Code: "TOTP_INVALID", Message: "TOTP 校验失败", HTTPStatus: http.StatusUnauthorized}
	ErrResetTokenInvalid  = &AppError{Code: "RESET_TOKEN_INVALID", Message: "重置令牌无效", HTTPStatus: http.StatusUnauthorized}
	ErrInternal           = &AppError{Code: "INTERNAL_ERROR", Message: "系统内部错误", HTTPStatus: http.StatusInternalServerError}
	ErrCircuitOpen        = &AppError{Code: "CIRCUIT_OPEN", Message: "系统繁忙，请稍后重试", HTTPStatus: http.StatusServiceUnavailable}
	ErrUserNotFoundByAuth = &AppError{Code: "USER_NOT_FOUND", Message: "用户不存在", HTTPStatus: http.StatusNotFound}
)

// AsAppError 将任意错误映射为可返回给客户端的 AppError
func AsAppError(err error) *AppError {
	if err == nil {
		return nil
	}
	var ae *AppError
	if errors.As(err, &ae) {
		return ae
	}
	return ErrInternal.WithErr(err)
}
