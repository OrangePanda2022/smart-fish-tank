package controller

import (
	"auth/internal/dto/request"
	"auth/internal/dto/response"
	"auth/internal/dto/result"
	errorHandler "auth/internal/error"
	"auth/internal/middleware"
	"auth/internal/service"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	hydra "github.com/ory/hydra-client-go/v25"
)

// Handler 聚合 HTTP 层依赖
type Handler struct {
	authSvc           service.AuthService
	userSvc           service.UserService
	hydraAdmin        hydra.OAuth2API
	clientSecrets     map[string]string
	redirectAllowlist map[string]map[string]struct{}
	// Session Cookie 有效期
	sessionCookieMaxAge int
	// Hydra 登录记忆时长
	hydraLoginRememberFor int64
	// Hydra 授权记忆时长
	hydraConsentRememberFor int64
}

// NewHandler 预处理客户端凭据与重定向白名单，避免请求路径重复做字符串清洗
func NewHandler(authSvc service.AuthService,
	userSvc service.UserService,
	hydraAdminURL string,
	clients map[string]string,
	redirectAllowlist map[string][]string,
	sessionCookieMaxAge int,
	hydraLoginRememberFor int64,
	hydraConsentRememberFor int64,
) *Handler {
	allowlist := make(map[string]map[string]struct{}, len(redirectAllowlist))

	for clientID, urls := range redirectAllowlist {
		set := make(map[string]struct{}, len(urls))
		for _, u := range urls {
			set[u] = struct{}{}
		}
		allowlist[clientID] = set
	}

	hydraClient := hydra.NewAPIClient(&hydra.Configuration{
		Servers: hydra.ServerConfigurations{{URL: strings.TrimRight(hydraAdminURL, "/")}},
	})

	return &Handler{
		authSvc:           authSvc,
		userSvc:           userSvc,
		hydraAdmin:        hydraClient.OAuth2API,
		clientSecrets:     clients,
		redirectAllowlist: allowlist,
	}
}

// Register 注册用户
func (h *Handler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(c, errorHandler.ErrBadRequest)
		return
	}
	res, err := h.authSvc.Register(c.Request.Context(), req.Email, req.Password, req.UserName, getRequestID(c))
	if err != nil {
		result.Error(c, err)
		return
	}
	result.Success(c, res)
}

// OAuth2LoginRequest 获取 Hydra login request，前端据此决定是否展示登录页
func (h *Handler) OAuth2LoginRequest(c *gin.Context) {
	challenge := strings.TrimSpace(c.Query("login_challenge"))
	if challenge == "" {
		result.Error(c, errorHandler.ErrBadRequest.WithErr(fmt.Errorf("missing login_challenge")))
		return
	}

	loginReq, _, err := h.hydraAdmin.GetOAuth2LoginRequest(c.Request.Context()).LoginChallenge(challenge).Execute()
	if err != nil {
		result.Error(c, errorHandler.ErrBadRequest.WithErr(fmt.Errorf("get hydra login request failed: %w", err)))
		return
	}

	// Hydra 标记 skip 且带 subject 时，直接 accept，减少登录页跳转
	if loginReq.GetSkip() && strings.TrimSpace(loginReq.GetSubject()) != "" {
		redirectTo, acceptErr := h.acceptHydraLogin(c, challenge, loginReq.GetSubject(), true, h.hydraLoginRememberFor)
		if acceptErr != nil {
			result.Error(c, acceptErr)
			return
		}
		result.Success(c, gin.H{"redirect_to": redirectTo, "skip": true})
		return
	}

	// 本地会话存在时，复用已登录用户直接完成 Hydra challenge
	if userID := getUserID(c); userID != "" {
		redirectTo, acceptErr := h.acceptHydraLogin(c, challenge, userID, true, h.hydraLoginRememberFor)
		if acceptErr != nil {
			result.Error(c, acceptErr)
			return
		}
		result.Success(c, gin.H{"redirect_to": redirectTo, "skip": true})
		return
	}
	client := loginReq.GetClient()
	result.Success(c, response.LoginRequestResult{
		LoginChallenge: challenge,
		Skip:           false,
		ClientID:       client.ClientId,
		RequestedScope: loginReq.GetRequestedScope(),
		RequestURL:     loginReq.GetRequestUrl(),
		OIDCContext:    loginReq.GetOidcContext(),
		RequestedAud:   loginReq.GetRequestedAccessTokenAudience(),
		LoginSessionID: loginReq.GetSessionId(),
	})
}

// Login 校验邮箱密码并提交 Hydra login accept
func (h *Handler) Login(c *gin.Context) {
	challenge := strings.TrimSpace(c.Query("login_challenge"))
	if challenge == "" {
		result.Error(c, errorHandler.ErrBadRequest.WithErr(fmt.Errorf("missing login_challenge")))
		return
	}

	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(c, errorHandler.ErrBadRequest)
		return
	}

	res, err := h.authSvc.Login(c.Request.Context(), req.Email, req.Password, getRequestID(c))
	if err != nil {
		result.Error(c, err)
		return
	}

	if strings.TrimSpace(res.SessionID) != "" {
		c.SetCookie(middleware.SessionCookieName, res.SessionID, h.sessionCookieMaxAge, "/", "", false, true)
	}

	redirectTo, acceptErr := h.acceptHydraLogin(c, challenge, res.Profile.UserID, true, h.hydraLoginRememberFor)
	if acceptErr != nil {
		result.Error(c, acceptErr)
		return
	}

	res.RedirectTo = redirectTo

	result.Success(c, res)
}

// OAuth2ConsentRequest 获取 consent challenge，并在可跳过场景自动 accept
func (h *Handler) OAuth2ConsentRequest(c *gin.Context) {
	challenge := strings.TrimSpace(c.Query("consent_challenge"))
	if challenge == "" {
		result.Error(c, errorHandler.ErrBadRequest.WithErr(fmt.Errorf("missing consent_challenge")))
		return
	}

	consentReq, _, err := h.hydraAdmin.GetOAuth2ConsentRequest(c.Request.Context()).ConsentChallenge(challenge).Execute()
	if err != nil {
		result.Error(c, errorHandler.ErrBadRequest.WithErr(fmt.Errorf("get hydra consent request failed: %w", err)))
		return
	}

	// 用户未登录时拒绝 consent，让 Hydra 走标准 OAuth 错误返回
	if strings.TrimSpace(consentReq.GetSubject()) == "" && getUserID(c) == "" {
		redirectTo, rejectErr := h.rejectHydraConsent(c, challenge, http.StatusUnauthorized, "login_required", "end-user must login before consent")
		if rejectErr != nil {
			result.Error(c, rejectErr)
			return
		}
		result.Success(c, gin.H{"redirect_to": redirectTo, "rejected": true})
		return
	}

	// skip=true 表示历史授权可复用，按最小权限直接 accept
	if consentReq.GetSkip() {
		profile, profileErr := h.loadProfileByConsentSubject(c, consentReq)
		if profileErr != nil {
			result.Error(c, profileErr)
			return
		}
		redirectTo, acceptErr := h.acceptHydraConsent(c, challenge, consentReq.GetRequestedScope(), consentReq.GetRequestedAccessTokenAudience(), true, h.hydraConsentRememberFor, profile)
		if acceptErr != nil {
			result.Error(c, acceptErr)
			return
		}
		result.Success(c, gin.H{"redirect_to": redirectTo, "skip": true})
		return
	}

	client := consentReq.GetClient()
	result.Success(c, response.ConsentRequestResult{
		ConsentChallenge: challenge,
		Skip:             false,
		Subject:          consentReq.GetSubject(),
		ClientID:         client.ClientId,
		ClientName:       client.ClientName,
		RequestedScope:   consentReq.GetRequestedScope(),
		RequestedAud:     consentReq.GetRequestedAccessTokenAudience(),
		RequestURL:       consentReq.GetRequestUrl(),
	})
}

// OAuth2ConsentSubmit 提交 consent 决策，转发到 Hydra accept/reject
func (h *Handler) OAuth2ConsentSubmit(c *gin.Context) {
	var req request.ConsentSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(c, errorHandler.ErrBadRequest)
		return
	}

	challenge := strings.TrimSpace(req.ConsentChallenge)
	consentReq, _, err := h.hydraAdmin.GetOAuth2ConsentRequest(c.Request.Context()).ConsentChallenge(challenge).Execute()
	if err != nil {
		fmt.Println(1)
		result.Error(c, errorHandler.ErrBadRequest.WithErr(fmt.Errorf("get hydra consent request failed: %w", err)))
		return
	}

	if !req.Accept {
		redirectTo, rejectErr := h.rejectHydraConsent(c, challenge, http.StatusForbidden, "access_denied", "end-user denied the consent request")
		if rejectErr != nil {
			result.Error(c, rejectErr)
			return
		}
		result.Success(c, gin.H{"redirect_to": redirectTo, "accepted": false})
		return
	}

	grantScope := chooseGrantedItems(consentReq.GetRequestedScope(), req.GrantScope)
	grantAudience := chooseGrantedItems(consentReq.GetRequestedAccessTokenAudience(), req.GrantAudience)

	if req.RememberFor < 0 {
		fmt.Println(2)
		result.Error(c, errorHandler.ErrBadRequest.WithErr(fmt.Errorf("remember_for must be >= 0")))
		return
	}
	rememberFor := req.RememberFor
	if req.Remember && rememberFor == 0 {
		rememberFor = h.hydraConsentRememberFor
	}

	profile, profileErr := h.loadProfileByConsentSubject(c, consentReq)
	if profileErr != nil {
		result.Error(c, profileErr)
		return
	}
	redirectTo, acceptErr := h.acceptHydraConsent(c, challenge, grantScope, grantAudience, req.Remember, rememberFor, profile)
	if acceptErr != nil {
		result.Error(c, acceptErr)
		return
	}
	result.Success(c, gin.H{"redirect_to": redirectTo, "accepted": true})
}

// Logout 注销当前会话
func (h *Handler) Logout(c *gin.Context) {
	var req request.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(c, errorHandler.ErrBadRequest)
		return
	}

	sid, _ := c.Cookie(middleware.SessionCookieName)
	if strings.TrimSpace(sid) != "" {
		if err := h.authSvc.Logout(c.Request.Context(), sid, getRequestID(c)); err != nil {
			result.Error(c, err)
			return
		}
	}
	c.SetCookie(middleware.SessionCookieName, "", -1, "/", "", false, true)
	result.Success(c, gin.H{"message": "logout success"})
}

// VerifyTOTP 作为找回密码链路的一步，校验成功后签发 reset_token
func (h *Handler) VerifyTOTP(c *gin.Context) {
	var req request.VerifyTOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(c, errorHandler.ErrBadRequest)
		return
	}
	resetToken, err := h.authSvc.VerifyTOTPForReset(c.Request.Context(), req.Email, req.Code, getRequestID(c))
	if err != nil {
		result.Error(c, err)
		return
	}
	result.Success(c, gin.H{"reset_token": resetToken})
}

// RequestResetPassword 发起密码重置流程，接口本身不暴露用户是否存在
func (h *Handler) RequestResetPassword(c *gin.Context) {
	var req request.ResetRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(c, errorHandler.ErrBadRequest)
		return
	}
	if err := h.authSvc.RequestResetPassword(c.Request.Context(), req.Email); err != nil {
		result.Error(c, err)
		return
	}
	result.Success(c, gin.H{"next_step": ""})
}

// ResetPassword 消费一次性 reset_token 并重置密码
func (h *Handler) ResetPassword(c *gin.Context) {
	var req request.ResetConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(c, errorHandler.ErrBadRequest)
		return
	}
	if err := h.authSvc.ResetPassword(c.Request.Context(), req.ResetToken, req.NewPassword, getRequestID(c)); err != nil {
		result.Error(c, err)
		return
	}
	result.Success(c, gin.H{"message": "password reset success"})
}

// GetProfile 返回当前登录用户资料
func (h *Handler) GetProfile(c *gin.Context) {
	userID, ok := mustUserID(c)
	if !ok {
		return
	}
	profile, err := h.userSvc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		result.Error(c, err)
		return
	}
	result.Success(c, profile)
}

// UpdateName 修改当前登录用户昵称
func (h *Handler) UpdateName(c *gin.Context) {
	userID, ok := mustUserID(c)
	if !ok {
		return
	}
	var req request.UpdateNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(c, errorHandler.ErrBadRequest)
		return
	}
	if err := h.userSvc.UpdateName(c.Request.Context(), userID, req.UserName); err != nil {
		result.Error(c, err)
		return
	}
	result.Success(c, gin.H{"message": "user name updated"})
}

// ChangePassword 修改当前登录用户密码，并使历史刷新令牌失效
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, ok := mustUserID(c)
	if !ok {
		return
	}
	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(c, errorHandler.ErrBadRequest)
		return
	}
	if err := h.userSvc.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword, getRequestID(c)); err != nil {
		result.Error(c, err)
		return
	}
	result.Success(c, gin.H{"message": "password updated"})
}

// DeleteUser 软删除当前账号
func (h *Handler) DeleteUser(c *gin.Context) {
	userID, ok := mustUserID(c)
	if !ok {
		return
	}
	if err := h.userSvc.DeleteUser(c.Request.Context(), userID); err != nil {
		result.Error(c, err)
		return
	}
	result.Success(c, gin.H{"message": "user deleted"})
}

// AdminSetUserStatus 管理员启用/禁用用户账号
func (h *Handler) AdminSetUserStatus(c *gin.Context) {
	operatorID := getUserID(c)
	targetUserID := c.Param("user_id")
	if operatorID == "" || targetUserID == "" {
		result.Error(c, errorHandler.ErrBadRequest)
		return
	}
	var req request.AdminSetUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(c, errorHandler.ErrBadRequest)
		return
	}
	if err := h.userSvc.SetUserStatus(c.Request.Context(), operatorID, targetUserID, req.Enabled); err != nil {
		result.Error(c, err)
		return
	}
	result.Success(c, gin.H{"message": "user status updated"})
}

// AdminSetUserRole 管理员授予或取消目标用户管理员权限
func (h *Handler) AdminSetUserRole(c *gin.Context) {
	operatorID := getUserID(c)
	targetUserID := c.Param("user_id")
	if operatorID == "" || targetUserID == "" {
		result.Error(c, errorHandler.ErrBadRequest)
		return
	}
	var req request.AdminSetUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Error(c, errorHandler.ErrBadRequest)
		return
	}
	if err := h.userSvc.SetUserRole(c.Request.Context(), operatorID, targetUserID, req.Role); err != nil {
		result.Error(c, err)
		return
	}
	result.Success(c, gin.H{"message": "user role updated"})
}

func getClientCredentials(c *gin.Context) (clientID string, clientSecret string) {
	clientID, clientSecret, ok := c.Request.BasicAuth()
	if ok {
		return clientID, clientSecret
	}
	return c.PostForm("client_id"), c.PostForm("client_secret")
}

func mustUserID(c *gin.Context) (string, bool) {
	userID := getUserID(c)
	if userID == "" {
		result.Error(c, errorHandler.ErrUnauthorized)
		return "", false
	}
	return userID, true
}

// getUserID 从上下文读取认证中间件写入的 user_id
func getUserID(c *gin.Context) string {
	uid, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := uid.(string)
	return userID
}

// getRequestID 从上下文读取请求链路 ID，便于贯穿日志与审计
func getRequestID(c *gin.Context) string {
	rid, _ := c.Get(middleware.ContextRequestIDKey)
	requestID, _ := rid.(string)
	return requestID
}

func (h *Handler) acceptHydraLogin(c *gin.Context, challenge string, subject string, remember bool, rememberFor int64) (string, error) {
	payload := hydra.NewAcceptOAuth2LoginRequest(strings.TrimSpace(subject))
	payload.SetRemember(remember)
	payload.SetRememberFor(rememberFor)

	redirect, _, err := h.hydraAdmin.
		AcceptOAuth2LoginRequest(c.Request.Context()).
		LoginChallenge(challenge).
		AcceptOAuth2LoginRequest(*payload).
		Execute()
	if err != nil {
		return "", errorHandler.ErrInternal.WithErr(fmt.Errorf("accept hydra login request failed: %w", err))
	}
	return redirect.GetRedirectTo(), nil
}

func (h *Handler) acceptHydraConsent(c *gin.Context, challenge string, grantScope []string, grantAudience []string, remember bool, rememberFor int64, profile any) (string, error) {
	payload := hydra.NewAcceptOAuth2ConsentRequest()
	payload.SetGrantScope(grantScope)
	payload.SetGrantAccessTokenAudience(grantAudience)
	payload.SetRemember(remember)
	payload.SetRememberFor(rememberFor)
	// 在 token session 注入必要用户信息，减少资源服务回源查询
	payload.SetSession(hydra.AcceptOAuth2ConsentRequestSession{
		AccessToken: profile,
	})

	redirect, _, err := h.hydraAdmin.
		AcceptOAuth2ConsentRequest(c.Request.Context()).
		ConsentChallenge(challenge).
		AcceptOAuth2ConsentRequest(*payload).
		Execute()
	if err != nil {
		return "", errorHandler.ErrInternal.WithErr(fmt.Errorf("accept hydra consent request failed: %w", err))
	}
	return redirect.GetRedirectTo(), nil
}

func (h *Handler) rejectHydraConsent(c *gin.Context, challenge string, statusCode int, oauthError string, description string) (string, error) {
	reject := hydra.NewRejectOAuth2Request()
	reject.SetStatusCode(int64(statusCode))
	reject.SetError(oauthError)
	reject.SetErrorDescription(description)

	redirect, _, err := h.hydraAdmin.
		RejectOAuth2ConsentRequest(c.Request.Context()).
		ConsentChallenge(challenge).
		RejectOAuth2Request(*reject).
		Execute()
	if err != nil {
		return "", errorHandler.ErrInternal.WithErr(fmt.Errorf("reject hydra consent request failed: %w", err))
	}
	return redirect.GetRedirectTo(), nil
}

func chooseGrantedItems(requested []string, granted []string) []string {
	// 仅允许授权请求中出现过的条目，防止前端越权提交 scope/audience
	if len(requested) == 0 {
		return nil
	}
	if len(granted) == 0 {
		return requested
	}
	allowed := make(map[string]struct{}, len(requested))
	for _, item := range requested {
		allowed[item] = struct{}{}
	}
	out := make([]string, 0, len(granted))
	for _, item := range granted {
		if _, ok := allowed[item]; ok {
			out = append(out, item)
		}
	}
	return out
}

func (h *Handler) loadProfileByConsentSubject(c *gin.Context, consentReq *hydra.OAuth2ConsentRequest) (gin.H, error) {
	subject := strings.TrimSpace(consentReq.GetSubject())
	if subject == "" {
		subject = getUserID(c)
	}
	if subject == "" {
		return nil, errorHandler.ErrUnauthorized.WithErr(fmt.Errorf("missing consent subject"))
	}

	profile, err := h.userSvc.GetProfile(c.Request.Context(), subject)
	if err != nil {
		return nil, err
	}
	// 输出到 id_token 的 claims 使用最小字段集合，避免泄漏内部模型细节
	return gin.H{
		"sub":   profile.UserID,
		"email": profile.UserEmail,
		"name":  profile.UserName,
		"role":  profile.UserRoles,
	}, nil
}
