package captcha

import (
	"crypto/rand"
	"math/big"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateCaptcha 生成验证码
func GenerateCaptcha() string {
	digits := "0123456789"
	code := make([]byte, 6)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			// 如果crypto/rand失败，使用时间戳作为fallback
			return string(code[:i]) + string(digits[i%10]) + string(code[i+1:])
		}
		code[i] = digits[n.Int64()]
	}
	return string(code)
}

// VerifyCaptcha 校验验证码
func VerifyCaptcha(code, captchaCode string) bool {
	return code == captchaCode
}

// GenerateInviteCode 生成邀请码
func GenerateInviteCode() string {
	code := make([]byte, 6)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// 如果crypto/rand失败，使用时间戳作为fallback
			return string(code[:i]) + string(charset[i%len(charset)]) + string(code[i+1:])
		}
		code[i] = charset[n.Int64()]
	}
	return string(code)
}
