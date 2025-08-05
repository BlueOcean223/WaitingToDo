package models

import "time"

type InviteCode struct {
	Id         int       `json:"id"`
	TaskId     int       `json:"task_id"`
	InviteCode string    `json:"invite_code"`
	CreateTime time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"autoUpdateTime"`
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
