package user

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type UserPerm struct {
	ID     string `gorm:"primaryKey;autoIncrement"`
	UserID string `gorm:"index"`
	Perm   string
}

func (u *UserPermsList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, u)
}

func (u UserPermsList) Value() (driver.Value, error) {
	return json.Marshal(u)
}

type UserPermsList []UserPerm
