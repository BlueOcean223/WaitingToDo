package models

type TaskNoticeHistory struct {
	Id     int `json:"id"`
	TaskId int `json:"task_id"`
}

// TableName 表名
func (TaskNoticeHistory) TableName() string {
	return "task_notice_history"
}
