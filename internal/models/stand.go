package models

const (
	StandTypeBuyer    string = "CONSUMPTION"
	StandTypeActivity string = "ACTIVITY"
)

type Stand struct {
	Id          int    `json:"id" db:"id"`
	UserId      int    `json:"user_id" db:"user_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Type        string `json:"type" db:"type"`
	Price       int    `json:"price" db:"price"`
	Stock       int    `json:"stock" db:"stock"`
}
