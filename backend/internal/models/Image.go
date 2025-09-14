package models

import "time"

type Image struct {
	ID         int       `json:"id" gorm:"type:int"`
	Md5        string    `json:"md5" gorm:"type:varchar(255)"`
	Url        string    `json:"url" gorm:"type:varchar(255)"`
	CreateTime time.Time `json:"create_time" gorm:"type:datetime;autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"type:datetime;autoUpdateTime"`
}
