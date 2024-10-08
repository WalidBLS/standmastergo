package models

const (
	TombolaStatusStarted = "STARTED"
	TombolaStatusEnded   = "ENDED"
)

type Tombola struct {
	Id         int    `json:"id" db:"id"`
	KermesseId int    `json:"kermesse_id" db:"kermesse_id"`
	Name       string `json:"name" db:"name"`
	Status     string `json:"status" db:"status"`
	Price      int    `json:"price" db:"price"`
	Gift       string `json:"gift" db:"gift"`
}
