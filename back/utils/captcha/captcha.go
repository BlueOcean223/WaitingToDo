package captcha

import "math/rand"

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
