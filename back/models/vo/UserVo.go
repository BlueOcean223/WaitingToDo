package vo

type UserVo struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Captcha     string `json:"captcha"`
	Description string `json:"description"`
	Pic         string `json:"pic"`
}
