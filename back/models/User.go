package models

import "time"

type User struct {
	Id          int       `json:"id"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	Name        string    `json:"name"`
	Pic         string    `json:"pic"`
	Description string    `json:"description"`
	CreateTime  time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime  time.Time `json:"update_time" gorm:"autoUpdateTime"`
}
