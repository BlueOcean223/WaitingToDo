package models

import "time"

type InviteCode struct {
	Id         int    `json:"id"`
	TaskId     int    `json:"task_id"`
	InviteCode string `json:"invite_code"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

func (ic *InviteCode) TableName() string {
	return "invite_codes"
}

func NewInviteCode(taskId int, inviteCode string) *InviteCode {
	return &InviteCode{
		TaskId:     taskId,
		InviteCode: inviteCode,
		CreateTime: time.Now().Format(time.DateTime),
		UpdateTime: time.Now().Format(time.DateTime),
	}
}
