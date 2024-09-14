package userConfig

type User struct {
	Login    string `json:"login" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required"`
}

type PutUser struct {
	Login    string `json:"login,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

type AuthData struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TableUser struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Date      string `json:"date"`
	IsBlocked bool   `json:"isBlocked"`
	IsAdmin   bool   `json:"isAdmin"`
}

type Meta struct {
	TotalAmount int    `json:"totalAmount"`
	SortBy      string `json:"sortBy"`
	SortOrder   string `json:"sortOrder"`
}

type MetaResponse struct {
	Data []TableUser `json:"data"`
	Meta Meta        `json:"meta"`
}

type GetAllQuery struct {
	SearchTerm string
	SortBy     string
	SortOrder  string
	IsBlocked  bool
	Limit      int
	Offset     int
}
