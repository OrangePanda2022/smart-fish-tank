package user

import (
	"auth/internal/util"
	"time"

	"gorm.io/gorm"
)

const (
	// RoleUser 表示普通用户角色
	RoleUser = "user"
	// RoleAdmin 表示管理员角色
	RoleAdmin = "admin"
)

// User 是用户聚合根，承载身份认证与用户资料的核心字段
type UserEntity struct {
	UserID           string        `gorm:"type:char(36);primaryKey"`
	UserEmail        string        `gorm:"size:255;uniqueIndex;not null"`
	UserPasswordHash string        `gorm:"size:255;not null"`
	UserName         string        `gorm:"size:100;not null"`
	UserRoles        UserRolesList `gorm:"type:TEXT"`
	UserPerms        UserPermsList `gorm:"type:TEXT"`
	UserStatus       string        `gorm:"size:16;not null;default:enabled"`
	TOTPSecret       string        `gorm:"size:128"`
	UserRevoked      time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

func (u *UserEntity) BeforeCreate(tx *gorm.DB) (err error) {
	// 如果 ID 已经有值，则不覆盖
	if u.UserID == "" {
		uid, err := util.NewUUIDv7()
		if err != nil {
			return err
		}
		u.UserID = uid
	}
	return nil
}
