package models

import "time"

type File struct {
	Id         int       `json:"id" gorm:"primaryKey"`
	TaskId     int       `json:"task_id" gorm:"index"`
	Name       string    `json:"name"`
	Url        string    `json:"url"`
	CreateTime time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"autoUpdateTime"`
}
