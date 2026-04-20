package util

import "github.com/google/uuid"

// NewUUIDv7 生成 UUID v7 字符串
func NewUUIDv7() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
