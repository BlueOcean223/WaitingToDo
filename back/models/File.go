package models

type File struct {
	Id         int    `json:"id"`
	TaskId     int    `json:"task_id"`
	Name       string `json:"name"`
	Url        string `json:"url"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}
