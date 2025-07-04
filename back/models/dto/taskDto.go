package dto

import "back/models"

type TaskDto struct {
	Id          int           `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Ddl         string        `json:"ddl"`
	Status      int           `json:"status"`
	Count       int64         `json:"count"`
	Attachments []models.File `json:"attachments"`
}
