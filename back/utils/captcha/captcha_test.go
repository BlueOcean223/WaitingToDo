package captcha

import "testing"

func TestGenerateCaptcha(t *testing.T) {
	captcha := GenerateCaptcha()
	t.Log(captcha)
}

func TestGenerateInviteCode(t *testing.T) {
	for i := 0; i < 10; i++ {
		inviteCode := GenerateInviteCode()
		t.Log(inviteCode)
	}
}
