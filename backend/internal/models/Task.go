package models

import "time"

type Task struct {
	Id          int       `json:"id" gorm:"primaryKey;type:int"`
	UserId      int       `json:"user_id" gorm:"index:idx_user_type_status,priority:1;index:idx_user_ddl;type:int"`
	Title       string    `json:"title" gorm:"type:varchar(255)"`
	Description string    `json:"description" gorm:"type:varchar(255)"`
	Ddl         string    `json:"ddl" gorm:"index:idx_ddl_status;index:idx_user_ddl,priority:2;type:varchar(255)"`
	Type        int       `json:"type" gorm:"index:idx_user_type_status,priority:2;type:int"`
	Status      int       `json:"status" gorm:"index:idx_user_type_status,priority:3;index:idx_ddl_status,priority:2;type:int"`
	CreateTime  time.Time `json:"create_time" gorm:"autoCreateTime;type:datetime"`
	UpdateTime  time.Time `json:"update_time" gorm:"autoUpdateTime;type:datetime"`
}
