package event

import "time"

// EventType 定义审计事件类型
type EventType string

const (
	EventLogin         EventType = "login"
	EventLogout        EventType = "logout"
	EventPasswordReset EventType = "password_reset"
	EventTOTPFailed    EventType = "totp_failed"
)

// AuditEvent 为审计日志统一消息结构
type AuditEvent struct {
	Type      EventType         `json:"type"`
	UserID    string            `json:"user_id,omitempty"`
	UserMail  string            `json:"user_email,omitempty"`
	Action    string            `json:"action,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
	Meta      map[string]string `json:"meta,omitempty"`
	At        time.Time         `json:"at"`
}
