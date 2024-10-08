package models

type contextKey string

const (
	UserIDKey   contextKey = "userId"
	UserRoleKey contextKey = "userRole"
)

const (
	UserRoleOrganizer   string = "ORGANIZER"
	UserRoleStandHolder string = "STAND_HOLDER"
	UserRoleParent      string = "PARENT"
	UserRoleChild       string = "CHILD"
)

type User struct {
	Id       int    `json:"id" db:"id"`
	ParentId *int   `json:"parentId" db:"parent_id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	Role     string `json:"role" db:"role"`
	Credit   int    `json:"credit" db:"credit"`
}

type UserBasic struct {
	Id     int    `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Email  string `json:"email" db:"email"`
	Role   string `json:"role" db:"role"`
	Credit int    `json:"credit" db:"credit"`
}

type UserMe struct {
	Id       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Role     string `json:"role" db:"role"`
	Credit   int    `json:"credit" db:"credit"`
	Token    string `json:"token"`
	HasStand bool   `json:"has_stand"`
}
