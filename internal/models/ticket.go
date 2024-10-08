package models

import "time"

type TicketUser struct {
	Id    int    `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	Role  string `json:"role" db:"role"`
}

type TicketTombola struct {
	Id     int    `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Status string `json:"status" db:"status"`
	Price  int    `json:"price" db:"price"`
	Gift   string `json:"gift" db:"gift"`
}

type TicketKermesse struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Status      string `json:"status" db:"status"`
}

type Ticket struct {
	Id        int            `json:"id" db:"id"`
	IsWinner  bool           `json:"is_winner" db:"is_winner"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	User      TicketUser     `json:"user" db:"user"`
	Tombola   TicketTombola  `json:"tombola" db:"tombola"`
	Kermesse  TicketKermesse `json:"kermesse" db:"kermesse"`
}
