package models

import "time"

type Image struct {
	ID         int       `json:"id"`
	Md5        string    `json:"md5"`
	Url        string    `json:"url"`
	CreateTime time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"autoUpdateTime"`
}
