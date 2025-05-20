package models

type Task struct {
	Id          int    `json:"id"`
	UserId      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Ddl         string `json:"ddl"`
	Type        int    `json:"type"`
	Status      int    `json:"status"`
}
