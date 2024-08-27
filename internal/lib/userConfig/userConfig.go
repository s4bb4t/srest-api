package userConfig

type User struct {
	Login    string `json:"login"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type AuthData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Login struct {
	Login string `json:"login"`
}

type TableUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Date     string `json:"date"`
	Block    bool   `json:"block"`
	Admin    bool   `json:"admin"`
}
