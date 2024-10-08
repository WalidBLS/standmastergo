package ticket

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"standmaster/internal/models"
)

type TicketRepository interface {
	FindAll(filters map[string]interface{}) ([]models.Ticket, error)
	FindById(id int) (models.Ticket, error)
	Create(input map[string]interface{}) error
	CanCreate(input map[string]interface{}) (bool, error)
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (s *Repository) FindAll(filters map[string]interface{}) ([]models.Ticket, error) {
	tickets := []models.Ticket{}
	query := `
		SELECT DISTINCT
			t.id AS id,
			t.is_winner AS is_winner,
			u.id AS "user.id",
			u.name AS "user.name",
			u.email AS "user.email",
			u.role AS "user.role",
			tb.id AS "tombola.id",
			tb.name AS "tombola.name",
			tb.status AS "tombola.status",
			tb.price AS "tombola.price",
			tb.gift AS "tombola.gift",
			k.id AS "kermesse.id",
			k.name AS "kermesse.name",
			k.description AS "kermesse.description",
			k.status AS "kermesse.status"
		FROM tickets t
		JOIN users u ON t.user_id = u.id
		JOIN tombolas tb ON t.tombola_id = tb.id
		JOIN kermesses k ON tb.kermesse_id = k.id
		WHERE 1=1
	`
	if filters["organizer_id"] != nil {
		query += fmt.Sprintf(" AND k.user_id IS NOT NULL AND k.user_id = %v", filters["organizer_id"])
	}
	if filters["parent_id"] != nil {
		query += fmt.Sprintf(" AND u.parent_id IS NOT NULL AND u.parent_id = %v", filters["parent_id"])
	}
	if filters["child_id"] != nil {
		query += fmt.Sprintf(" AND t.user_id IS NOT NULL AND t.user_id = %v", filters["child_id"])
	}
	err := s.db.Select(&tickets, query)

	return tickets, err
}

func (s *Repository) FindById(id int) (models.Ticket, error) {
	ticket := models.Ticket{}
	query := `
		SELECT
			t.id AS id,
			t.is_winner AS is_winner,
			u.id AS "user.id",
			u.name AS "user.name",
			u.email AS "user.email",
			u.role AS "user.role",
			tb.id AS "tombola.id",
			tb.name AS "tombola.name",
			tb.status AS "tombola.status",
			tb.price AS "tombola.price",
			tb.gift AS "tombola.gift",
			k.id AS "kermesse.id",
			k.name AS "kermesse.name",
			k.description AS "kermesse.description",
			k.status AS "kermesse.status"
		FROM tickets t
		JOIN users u ON t.user_id = u.id
		JOIN tombolas tb ON t.tombola_id = tb.id
		JOIN kermesses k ON tb.kermesse_id = k.id
		WHERE t.id=$1
	`
	err := s.db.Get(&ticket, query, id)

	return ticket, err
}

func (s *Repository) CanCreate(input map[string]interface{}) (bool, error) {
	var isAssociated bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM kermesses_users ku
			JOIN kermesses k ON k.id = ku.kermesse_id
			WHERE ku.kermesse_id = $1 AND ku.user_id = $2 AND k.status = $3
		) AS is_associated
	`
	err := s.db.QueryRow(query, input["kermesse_id"], input["user_id"], models.KermesseStatusStarted).Scan(&isAssociated)

	return isAssociated, err
}

func (s *Repository) Create(input map[string]interface{}) error {
	query := "INSERT INTO tickets (user_id, tombola_id) VALUES ($1, $2)"
	_, err := s.db.Exec(query, input["user_id"], input["tombola_id"])

	return err
}
