package models

import "time"

type Friend struct {
	Id         int       `json:"id"`
	UserId     int       `json:"user_id"`
	FriendId   int       `json:"friend_id"`
	Status     int       `json:"status"`
	CreateTime time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime time.Time `json:"update_time" gorm:"autoUpdateTime"`
}
