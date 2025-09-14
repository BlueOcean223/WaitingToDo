package models

import "time"

type Message struct {
	Id          int       `json:"id" gorm:"primaryKey;type:int"`
	Title       string    `json:"title" gorm:"type:varchar(255)"`
	Description string    `json:"description" gorm:"type:varchar(255)"`
	FromId      int       `json:"from_id" gorm:"index;type:int"`
	ToId        int       `json:"to_id" gorm:"index:idx_to_read;index:idx_to_time;type:int"`
	Type        int       `json:"type" gorm:"type:int"`
	SendTime    string    `json:"send_time" gorm:"index:idx_to_time,priority:2;type:varchar(255)"`
	OutId       int       `json:"out_id" gorm:"type:int"`
	IsRead      int       `json:"is_read" gorm:"index:idx_to_read,priority:2;type:int"`
	CreateTime  time.Time `json:"create_time" gorm:"autoCreateTime;type:datetime"`
	UpdateTime  time.Time `json:"update_time" gorm:"autoUpdateTime;type:datetime"`
}
