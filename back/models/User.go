package models

type User struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	Pic         string `json:"pic"`
	Description string `json:"description"`
}
