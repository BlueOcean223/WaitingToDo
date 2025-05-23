package models

type Friend struct {
	Id       int `json:"id"`
	UserId   int `json:"user_id"`
	FriendId int `json:"friend_id"`
	Status   int `json:"status"`
}
