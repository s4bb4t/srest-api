package userConfig

type User struct {
	Login    string `json:"login" validate:"required,min=5,max=30,alphanum,excludesall= "`
	Username string `json:"username" validate:"required,min=5,max=30,excludesall= "`
	Password string `json:"password" validate:"required,min=3,max=30,excludesall= "`
	Email    string `json:"email" validate:"required,min=5,max=30,excludesall= ~!#$%^&*()-_=+[{]};:'<>/"`
}

type AuthData struct {
	Login    string `json:"login" validate:"required,min=5,max=30,alphanum,excludesall= "`
	Password string `json:"password" validate:"required,min=5,max=30,excludesall= "`
}

type TableUser struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Date      string `json:"date"`
	IsBlocked bool   `json:"isBlocked"`
	IsAdmin   bool   `json:"isAdmin"`
}
