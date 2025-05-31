package dto

type TeamUserDto struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Pic    string `json:"pic"`
	Status int    `json:"status"`
}

type TeamTaskDto struct {
	Id          int           `json:"id"`
	UserId      int           `json:"user_id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Ddl         string        `json:"ddl"`
	Status      int           `json:"status"`
	Users       []TeamUserDto `json:"users"`
}
