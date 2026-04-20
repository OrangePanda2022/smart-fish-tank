package security

import "golang.org/x/crypto/bcrypt"

// HashPassword 使用 bcrypt 对密码进行加盐哈希
func HashPassword(plain string) (string, error) {
	data, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// CheckPassword 校验明文密码与哈希值是否匹配
func CheckPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}
