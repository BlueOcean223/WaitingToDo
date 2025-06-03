package hashPassword

import "testing"

func TestHashPassword(t *testing.T) {
	// 测试加密
	password := "123456"
	hash, err := HashPassword(password)
	if err != nil {
		t.Error(err)
	}
	t.Log(hash)
	// 测试校验
	if CheckPasswordHash(password, hash) {
		t.Log("密码匹配")
	} else {
		t.Error("密码不匹配")
	}
}
