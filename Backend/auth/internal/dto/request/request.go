package request

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	UserName string `json:"user_name" binding:"required,min=2,max=50"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

// RefreshRequest 刷新令牌请求
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutRequest 登出请求，需传递 refresh token 以拉黑 JTI
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// VerifyTOTPRequest TOTP 校验请求
type VerifyTOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6,numeric"`
}

// UpdateNameRequest 修改昵称请求
type UpdateNameRequest struct {
	UserName string `json:"user_name" binding:"required,min=2,max=50"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=8,max=72"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=72"`
}

// ResetRequestRequest 发起重置密码请求
type ResetRequestRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetConfirmRequest 确认重置密码请求
type ResetConfirmRequest struct {
	ResetToken  string `json:"reset_token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=72"`
}

// AdminSetUserStatusRequest 管理员启用/禁用用户请求
type AdminSetUserStatusRequest struct {
	Enabled bool `json:"enabled"`
}

// AdminSetUserRoleRequest 管理员调整用户角色请求
type AdminSetUserRoleRequest struct {
	Role []string `json:"role" binding:"required"`
}

type ConsentSubmitRequest struct {
	ConsentChallenge string   `json:"consent_challenge" binding:"required"`
	Accept           bool     `json:"accept"`
	GrantScope       []string `json:"grant_scope"`
	GrantAudience    []string `json:"grant_audience"`
	Remember         bool     `json:"remember"`
	RememberFor      int64    `json:"remember_for"`
}
