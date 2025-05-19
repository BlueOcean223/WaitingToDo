package models

type Image struct {
	ID  int    `json:"id"`
	Md5 string `json:"md5"`
	Url string `json:"url"`
}
