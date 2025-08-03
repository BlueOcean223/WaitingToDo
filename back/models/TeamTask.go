package models

type TeamTask struct {
	Id     int `json:"id" gorm:"primaryKey"`
	TaskId int `json:"task_id" gorm:"index:idx_task_user;index"`
	UserId int `json:"user_id" gorm:"index:idx_task_user,priority:2;index"`
	Status int `json:"status"`
}

func (TeamTask) TableName() string {
	return "team_task"
}
