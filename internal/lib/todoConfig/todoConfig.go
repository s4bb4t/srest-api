package todoconfig

type Todo struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Created string `json:"created"`
	IsDone  bool   `json:"isdone"`
}

type Todos []Todo

type TodoRequest struct {
	Title string `json:"title"`
}
