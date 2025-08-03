package models

type TaskNoticeHistory struct {
	Id         int    `json:"id"`
	TaskId     int    `json:"task_id"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

// TableName 表名
func (TaskNoticeHistory) TableName() string {
	return "task_notice_history"
}
