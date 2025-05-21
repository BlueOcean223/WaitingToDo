package dto

type UserDto struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Pic         string `json:"pic"`
	Description string `json:"description"`
}
