package user

import "time"

// UserProfile 返回给前端的用户资料
type UserProfile struct {
	UserID     string    `json:"user_id"`
	UserEmail  string    `json:"user_email"`
	UserName   string    `json:"user_name"`
	UserRoles  []string  `json:"role"`
	UserStatus string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
