package models

type File struct {
	Id         int    `json:"id" gorm:"primaryKey"`
	TaskId     int    `json:"task_id" gorm:"index"`
	Name       string `json:"name"`
	Url        string `json:"url"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}
