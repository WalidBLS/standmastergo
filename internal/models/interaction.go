package models

import "time"

const (
	InteractionTypeConsumption string = "CONSUMPTION"
	InteractionTypeActivity    string = "ACTIVITY"

	InteractionStatusStarted string = "STARTED"
	InteractionStatusEnded   string = "ENDED"
)

type InteractionUser struct {
	Id    int    `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	Role  string `json:"role" db:"role"`
}

type InteractionStand struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Type        string `json:"type" db:"type"`
	Price       int    `json:"price" db:"price"`
}

type InteractionKermesse struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Status      string `json:"status" db:"status"`
}

type Interaction struct {
	Id        int                 `json:"id" db:"id"`
	Type      string              `json:"type" db:"type"`
	Status    string              `json:"status" db:"status"`
	Credit    int                 `json:"credit" db:"credit"`
	Point     int                 `json:"point" db:"point"`
	CreatedAt time.Time           `json:"created_at" db:"created_at"`
	User      InteractionUser     `json:"user" db:"user"`
	Stand     InteractionStand    `json:"stand" db:"stand"`
	Kermesse  InteractionKermesse `json:"kermesse" db:"kermesse"`
}

type InteractionBasic struct {
	Id        int              `json:"id" db:"id"`
	Type      string           `json:"type" db:"type"`
	Status    string           `json:"status" db:"status"`
	Credit    int              `json:"credit" db:"credit"`
	Point     int              `json:"point" db:"point"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
	User      InteractionUser  `json:"user" db:"user"`
	Stand     InteractionStand `json:"stand" db:"stand"`
}
