package models

import "time"

type User struct {
	Id          int       `json:"id" gorm:"type:int"`
	Email       string    `json:"email" gorm:"type:varchar(255)"`
	Password    string    `json:"password" gorm:"type:varchar(255)"`
	Name        string    `json:"name" gorm:"type:varchar(255)"`
	Pic         string    `json:"pic" gorm:"type:varchar(255)"`
	Description string    `json:"description" gorm:"type:varchar(255)"`
	CreateTime  time.Time `json:"create_time" gorm:"type:datetime;autoCreateTime"`
	UpdateTime  time.Time `json:"update_time" gorm:"type:datetime;autoUpdateTime"`
}
