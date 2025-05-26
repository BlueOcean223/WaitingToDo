package models

type QueueActionType int

const (
	FriendRequestReject QueueActionType = iota
	FriendRequestAccept
)

type MQMessage struct {
	MessageType int             `json:"message_type"` // 对应Message.Type
	ActionType  QueueActionType `json:"action_type"`  // 操作类型
	RelationID  int             `json:"relation_id"`  // 对应Message.OutId
	RequesterID int             `json:"requester_id"` // 对应Message.FromId
	ReceiverID  int             `json:"receiver_id"`  // 对应Message.ToId
}
