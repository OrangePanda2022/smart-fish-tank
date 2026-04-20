package user

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type UserRole struct {
	ID     string `gorm:"primaryKey;autoIncrement"`
	UserID string `gorm:"index"`
	Role   string
}

func (u *UserRolesList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, u)
}

func (u UserRolesList) Value() (driver.Value, error) {
	return json.Marshal(u)
}

type UserRolesList []UserRole
