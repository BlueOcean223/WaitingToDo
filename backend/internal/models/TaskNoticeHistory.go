package models

import "time"

type TaskNoticeHistory struct {
	Id         int       `json:"id" gorm:"type:int"`
	TaskId     int       `json:"task_id" gorm:"type:int"`
	CreateTime time.Time `json:"create_time" gorm:"type:datetime;autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"type:datetime;autoUpdateTime"`
}

// TableName 表名
func (TaskNoticeHistory) TableName() string {
	return "task_notice_history"
}
