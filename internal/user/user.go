package user

type User struct {
	AuthData
	Email string `json:"email"`
}

type AuthData struct {
	Login
	Password string `json:"password"`
}

type Login struct {
	Username string `json:"username"`
}

type TableUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Date     string `json:"date"`
	Blocked  bool   `json:"blocked"`
	Admin    bool   `json:"admin"`
}
