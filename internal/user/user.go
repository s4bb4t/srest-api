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
