package models

const (
	KermesseStatusStarted string = "STARTED"
	KermesseStatusEnded   string = "ENDED"
)

type Kermesse struct {
	Id          int    `json:"id" db:"id"`
	UserId      int    `json:"user_id" db:"user_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Status      string `json:"status" db:"status"`
}

type KermesseStats struct {
	StandCount        int `json:"stand_count"`
	TombolaCount      int `json:"tombola_count"`
	UserCount         int `json:"user_count"`
	InteractionCount  int `json:"interaction_count"`
	InteractionIncome int `json:"interaction_income"`
	TombolaIncome     int `json:"tombola_income"`
	Points            int `json:"points"`
}

type KermesseWithStats struct {
	Id                int    `json:"id" db:"id"`
	UserId            int    `json:"user_id" db:"user_id"`
	Name              string `json:"name" db:"name"`
	Description       string `json:"description" db:"description"`
	Status            string `json:"status" db:"status"`
	StandCount        int    `json:"stand_count"`
	TombolaCount      int    `json:"tombola_count"`
	UserCount         int    `json:"user_count"`
	InteractionCount  int    `json:"interaction_count"`
	InteractionIncome int    `json:"interaction_income"`
	TombolaIncome     int    `json:"tombola_income"`
	Points            int    `json:"points"`
}
