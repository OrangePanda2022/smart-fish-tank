package util

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

/*
主要修改说明（不改变业务行为）：
1. 函数拆分/合并：提取 TOTP 参数常量，统一生成与校验配置
2. 包结构调整：未调整包结构
3. 校验逻辑增删：保留原有容错窗口和算法配置
4. 注释修正：修正注释表述为实际行为
5. 代码风格统一：统一命名、常量与构造顺序
*/

const (
	totpPeriodSeconds = 30
	totpSkewWindow    = 1
	totpDigits        = otp.DigitsSix
	totpAlgorithm     = otp.AlgorithmSHA1
	totpSecretSize    = 20
)

// NewTOTPSecret 生成 TOTP 密钥和标准 otpauth URL
func NewTOTPSecret(issuer, email string) (string, string, error) {
	opts := totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: email,
		Algorithm:   totpAlgorithm,
		Digits:      totpDigits,
		Period:      totpPeriodSeconds,
		SecretSize:  totpSecretSize,
	}
	key, err := totp.Generate(opts)
	if err != nil {
		return "", "", fmt.Errorf("generate totp secret failed: %w", err)
	}

	// Secret 统一转大写，便于跨端传输和输入
	return strings.ToUpper(key.Secret()), key.URL(), nil
}

// VerifyTOTP 校验动态口令，允许前后一个周期的时间偏移
func VerifyTOTP(secret, code string) bool {
	ok, err := totp.ValidateCustom(code, secret, time.Now().UTC(), totp.ValidateOpts{
		Period:    totpPeriodSeconds,
		Skew:      totpSkewWindow,
		Digits:    totpDigits,
		Algorithm: totpAlgorithm,
	})
	return err == nil && ok
}

// BuildTOTPURL 生成兼容主流 Authenticator 的 otpauth URL
func BuildTOTPURL(secret, issuer, email string) string {
	label := url.PathEscape(fmt.Sprintf("%s:%s", issuer, email))

	q := url.Values{}
	q.Set("secret", secret)
	q.Set("issuer", issuer)
	q.Set("algorithm", "SHA1")
	q.Set("digits", fmt.Sprintf("%d", totpDigits))
	q.Set("period", fmt.Sprintf("%d", totpPeriodSeconds))

	return fmt.Sprintf("otpauth://totp/%s?%s", label, q.Encode())
}
