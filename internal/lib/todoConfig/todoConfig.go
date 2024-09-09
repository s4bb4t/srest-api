package todoconfig

type Todo struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Created string `json:"created"`
	IsDone  bool   `json:"isdone"`
}

type Todos []Todo

type TodoRequest struct {
	Title  string `json:"title"`
	IsDone bool   `json:"isdone"`
}

type TodoInfo struct {
	All       int `json:"all"`
	Completed int `json:"completed"`
	InWork    int `json:"inwork"`
}

type Meta struct {
	TotalAmount int `json:"total_amount"`
}

type MetaResponse struct {
	Data []Todo   `json:"data"`
	Info TodoInfo `json:"info"`
	Meta Meta     `json:"meta"`
}
