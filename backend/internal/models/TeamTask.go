package models

import "time"

type TeamTask struct {
	Id         int       `json:"id" gorm:"primaryKey;type:int"`
	TaskId     int       `json:"task_id" gorm:"index:idx_task_user;index;type:int"`
	UserId     int       `json:"user_id" gorm:"index:idx_task_user,priority:2;index;type:int"`
	Status     int       `json:"status" gorm:"type:int"`
	CreateTime time.Time `json:"create_time" gorm:"autoCreateTime;type:datetime"`
	UpdateTime time.Time `json:"update_time" gorm:"autoUpdateTime;type:datetime"`
}

func (TeamTask) TableName() string {
	return "team_task"
}
