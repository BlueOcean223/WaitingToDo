package models

import "time"

type TaskNoticeHistory struct {
	Id         int       `json:"id"`
	TaskId     int       `json:"task_id"`
	CreateTime time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"autoUpdateTime"`
}

// TableName 表名
func (TaskNoticeHistory) TableName() string {
	return "task_notice_history"
}
