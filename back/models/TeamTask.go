package models

type TeamTask struct {
	Id     int `json:"id"`
	TaskId int `json:"task_id"`
	UserId int `json:"user_id"`
	Status int `json:"status"`
}

func (TeamTask) TableName() string {
	return "team_task"
}
