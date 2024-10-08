package interaction

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"standmaster/internal/models"
)

type InteractionRepository interface {
	FindAll(filters map[string]interface{}) ([]models.InteractionBasic, error)
	FindById(id int) (models.Interaction, error)
	CanCreate(input map[string]interface{}) (bool, error)
	Create(input map[string]interface{}) error
	Update(id int, input map[string]interface{}) error
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (s *Repository) FindAll(filters map[string]interface{}) ([]models.InteractionBasic, error) {
	interactions := []models.InteractionBasic{}
	query := `
		SELECT DISTINCT
			i.id AS id,
			i.type AS type,
			i.status AS status,
			i.credit AS credit,
			i.point AS point,
			i.created_at AS created_at,
			u.id AS "user.id",
			u.name AS "user.name",
			u.email AS "user.email",
			u.role AS "user.role",
			s.id AS "stand.id",
			s.name AS "stand.name",
			s.description AS "stand.description",
			s.type AS "stand.type",
			s.price AS "stand.price"
		FROM interactions i
		JOIN users u ON i.user_id = u.id
		JOIN stands s ON i.stand_id = s.id
		WHERE 1=1
	`
	if filters["kermesse_id"] != nil {
		query += fmt.Sprintf(" AND i.kermesse_id = %v", filters["kermesse_id"])
	}
	if filters["parent_id"] != nil {
		query += fmt.Sprintf(" AND (u.id = %v OR u.parent_id = %v)", filters["parent_id"], filters["parent_id"])
	}
	if filters["child_id"] != nil {
		query += fmt.Sprintf(" AND u.id = %v", filters["child_id"])
	}
	if filters["stand_holder_id"] != nil {
		query += fmt.Sprintf(" AND s.user_id = %v", filters["stand_holder_id"])
	}
	query += " ORDER BY i.created_at DESC"
	err := s.db.Select(&interactions, query)

	return interactions, err
}

func (s *Repository) FindById(id int) (models.Interaction, error) {
	interaction := models.Interaction{}
	query := `
		SELECT
			i.id AS id,
			i.type AS type,
			i.status AS status,
			i.credit AS credit,
			i.point AS point,
			i.created_at AS created_at,
			u.id AS "user.id",
			u.name AS "user.name",
			u.email AS "user.email",
			u.role AS "user.role",
			s.id AS "stand.id",
			s.name AS "stand.name",
			s.description AS "stand.description",
			s.type AS "stand.type",
			s.price AS "stand.price",
			k.id AS "kermesse.id",
			k.name AS "kermesse.name",
			k.description AS "kermesse.description",
			k.status AS "kermesse.status"
		FROM interactions i
		JOIN users u ON i.user_id = u.id
		JOIN stands s ON i.stand_id = s.id
		JOIN kermesses k ON i.kermesse_id = k.id
		WHERE i.id=$1
	`
	err := s.db.Get(&interaction, query, id)

	return interaction, err
}

func (s *Repository) CanCreate(input map[string]interface{}) (bool, error) {
	var isAssociated bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM kermesses_users ku
  		JOIN kermesses_stands ks ON ku.kermesse_id = ks.kermesse_id
			JOIN kermesses k ON ku.kermesse_id = k.id
  		WHERE ku.user_id = $1 AND ks.stand_id = $2 AND k.status = $3
		) AS is_associated
 	`
	err := s.db.QueryRow(query, input["user_id"], input["stand_id"], models.KermesseStatusStarted).Scan(&isAssociated)

	return isAssociated, err
}

func (s *Repository) Create(input map[string]interface{}) error {
	query := "INSERT INTO interactions (user_id, kermesse_id, stand_id, type, credit) VALUES ($1, $2, $3, $4, $5)"
	_, err := s.db.Exec(query, input["user_id"], input["kermesse_id"], input["stand_id"], input["type"], input["credit"])

	return err
}

func (s *Repository) Update(id int, input map[string]interface{}) error {
	query := "UPDATE interactions SET status=$1, point=$2 WHERE id=$3"
	_, err := s.db.Exec(query, input["status"], input["point"], id)

	return err
}
