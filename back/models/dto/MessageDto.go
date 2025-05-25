package dto

type MessageDto struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	FromId      int    `json:"from_id"`
	ToId        int    `json:"to_id"`
	Type        int    `json:"type"` // 0普通消息，1好友请求，2组队邀请
	SendTime    string `json:"send_time"`
	OutId       int    `json:"out_id"`
	IsRead      int    `json:"is_read"` //0为未读，1为已读
}
