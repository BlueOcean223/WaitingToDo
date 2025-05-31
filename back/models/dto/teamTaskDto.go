package dto

type TeamTaskDto struct {
	Id          int       `json:"id"`
	UserId      int       `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Ddl         string    `json:"ddl"`
	Status      int       `json:"status"`
	Users       []UserDto `json:"users"`
}
