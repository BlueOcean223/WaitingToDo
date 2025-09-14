package models

import "time"

type InviteCode struct {
	Id         int       `json:"id" gorm:"type:int"`
	TaskId     int       `json:"task_id" gorm:"type:int"`
	InviteCode string    `json:"invite_code" gorm:"type:varchar(255)"`
	CreateTime time.Time `json:"create_time" gorm:"type:datetime;autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"type:datetime;autoUpdateTime"`
}

func (ic *InviteCode) TableName() string {
	return "invite_codes"
}

func NewInviteCode(taskId int, inviteCode string) *InviteCode {
	return &InviteCode{
		TaskId:     taskId,
		InviteCode: inviteCode,
	}
}
