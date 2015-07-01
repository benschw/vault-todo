package todo

type Todo struct {
	Id      int    `json:"id"`
	Status  string `json:"status"`
	Content string `json:"content"`
}
