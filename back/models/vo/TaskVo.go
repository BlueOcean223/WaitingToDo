package vo

type TaskVo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Ddl         string `json:"ddl"`
	Type        int    `json:"type"`
}
