package captcha

import "testing"

func TestGenerateCaptcha(t *testing.T) {
	captcha := GenerateCaptcha()
	t.Log(captcha)
}
