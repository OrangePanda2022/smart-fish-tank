package response

import (
	"auth/internal/domain/user"

	hydra "github.com/ory/hydra-client-go/v25"
)

// RegisterResult 为注册接口输出
type RegisterResult struct {
	Profile   user.UserProfile `json:"profile"`
	SessionID string           `json:"session_id"`
}

// LoginResult 为登录接口输出
type LoginResult struct {
	Profile    user.UserProfile `json:"profile"`
	SessionID  string           `json:"session_id"`
	RedirectTo string           `json:"redirect_to"`
}

type TOTPResult struct {
	TOTPURL string `json:"totp_url"`
}

type LoginRequestResult struct {
	LoginChallenge string                                         `json:"login_challenge"`
	Skip           bool                                           `json:"skip"`
	ClientID       *string                                        `json:"client_id"`
	RequestedScope []string                                       `json:"requested_scope"`
	RequestURL     string                                         `json:"request_url"`
	OIDCContext    hydra.OAuth2ConsentRequestOpenIDConnectContext `json:"oidc_context"`
	RequestedAud   []string                                       `json:"requested_aud"`
	LoginSessionID string                                         `json:"login_session_id"`
}

type ConsentRequestResult struct {
	ConsentChallenge string   `json:"consent_challenge"`
	Skip             bool     `json:"skip"`
	Subject          string   `json:"subject"`
	ClientID         *string  `json:"client_id"`
	ClientName       *string  `json:"client_name"`
	RequestedScope   []string `json:"requested_scope"`
	RequestedAud     []string `json:"requested_aud"`
	RequestURL       string   `json:"request_url"`
}
