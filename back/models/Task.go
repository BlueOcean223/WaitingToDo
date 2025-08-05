package models

import "time"

type Task struct {
	Id          int       `json:"id" gorm:"primaryKey"`
	UserId      int       `json:"user_id" gorm:"index:idx_user_type_status,priority:1;index:idx_user_ddl"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Ddl         string    `json:"ddl" gorm:"index:idx_ddl_status;index:idx_user_ddl,priority:2"`
	Type        int       `json:"type" gorm:"index:idx_user_type_status,priority:2"`
	Status      int       `json:"status" gorm:"index:idx_user_type_status,priority:3;index:idx_ddl_status,priority:2"`
	CreateTime  time.Time `json:"create_time" gorm:"autoCreateTime"`
	UpdateTime  time.Time `json:"update_time" gorm:"autoUpdateTime"`
}
