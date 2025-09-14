package models

import "time"

type Friend struct {
	Id         int       `json:"id" gorm:"type:int"`
	UserId     int       `json:"user_id" gorm:"index:idx_user_friend_status,priority:1;type:int"`
	FriendId   int       `json:"friend_id" gorm:"index:idx_user_friend_status,priority:2;type:int"`
	Status     int       `json:"status" gorm:"index:idx_user_friend_status,priority:3;type:int"`
	CreateTime time.Time `json:"create_time" gorm:"autoCreateTime;type:datetime"`
	UpdateTime time.Time `json:"update_time" gorm:"autoUpdateTime;type:datetime"`
}
