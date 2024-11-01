package userConfig

type User struct {
	Login       string `json:"login" validate:"required,min=2,max=60,alpha"`
	Username    string `json:"username" validate:"required,min=1,max=60,alphanumunicode"`
	Password    string `json:"password" validate:"required,min=6,max=60,alphanumunicode"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phoneNumber" validate:"omitempty,e164"`
}

type PutUser struct {
	Username    string `json:"username,omitempty" validate:"omitempty,min=1,max=60,alphanumunicode"`
	Email       string `json:"email,omitempty" validate:"omitempty,email"`
	PhoneNumber string `json:"phoneNumber,omitempty" validate:"omitempty,e164"`
}

type Pwd struct {
	Password string `json:"password" validate:"required,min=6,max=60,alphanumunicode"`
}

type AuthData struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TableUser struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Date        string `json:"date"`
	IsBlocked   bool   `json:"isBlocked"`
	IsAdmin     bool   `json:"isAdmin"`
	PhoneNumber string `json:"phoneNumber"`
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
	IsBlocked  *bool
	Limit      int
	Offset     int
}
