package dto

type UserDto struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	Pic         string `json:"pic"`
	Description string `json:"description"`
}
