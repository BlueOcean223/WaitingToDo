package captcha

import "math/rand"

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateCaptcha 生成验证码
func GenerateCaptcha() string {
	digits := "0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = digits[rand.Intn(len(digits))]
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
		idx := rand.Intn(len(charset))
		code[i] = charset[idx]
	}
	return string(code)
}
