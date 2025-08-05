package models

import "time"

type Message struct {
	Id          int       `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	FromId      int       `json:"from_id" gorm:"index"`
	ToId        int       `json:"to_id" gorm:"index:idx_to_read;index:idx_to_time"`
	Type        int       `json:"type"`
	SendTime    string    `json:"send_time" gorm:"index:idx_to_time,priority:2"`
	OutId       int       `json:"out_id"`
	IsRead      int       `json:"is_read" gorm:"index:idx_to_read,priority:2"`
	CreateTime  time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime  time.Time `json:"update_time" gorm:"autoUpdateTime"`
}
