package models

import "time"

type File struct {
	Id         int       `json:"id" gorm:"primaryKey;type:int"`
	TaskId     int       `json:"task_id" gorm:"type:int"`
	Name       string    `json:"name" gorm:"type:varchar(255)"`
	Url        string    `json:"url" gorm:"type:varchar(255)"`
	CreateTime time.Time `json:"create_time" gorm:"type:datetime;autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"type:datetime;autoUpdateTime"`
}
