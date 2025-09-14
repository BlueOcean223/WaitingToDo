package dto

type FriendInfoDto struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Pic         string `json:"pic"`
	Description string `json:"description"`
	IsFriend    int    `json:"isFriend"` // 0为待同意，1为已是好友，2为未添加
}
