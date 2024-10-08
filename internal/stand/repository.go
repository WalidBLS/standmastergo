package stand

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"standmaster/internal/models"
)

type StandRepository interface {
	FindAll(filters map[string]interface{}) ([]models.Stand, error)
	FindById(id int) (models.Stand, error)
	FindByUserId(id int) (models.Stand, error)
	Create(input map[string]interface{}) error
	Update(id int, input map[string]interface{}) error
	UpdateByUserId(userId int, input map[string]interface{}) error
	UpdateStock(id int, n int) error
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (s *Repository) FindAll(filters map[string]interface{}) ([]models.Stand, error) {
	stands := []models.Stand{}
	query := `
		SELECT DISTINCT
			s.id AS id,
			s.user_id AS user_id,
			s.name AS name,
			s.description AS description,
			s.type AS type,
			s.price AS price,
			s.stock AS stock
		FROM stands s
		LEFT JOIN kermesses_stands ks ON s.id = ks.stand_id
		WHERE 1=1 AND s.id IS NOT NULL
	`
	if filters["kermesse_id"] != nil {
		query += fmt.Sprintf(" AND ks.kermesse_id IS NOT NULL AND ks.kermesse_id = %v", filters["kermesse_id"])
	}
	if filters["is_free"] != nil {
		query += `
			AND (
				ks.kermesse_id IS NULL
				OR s.id NOT IN (
					SELECT ks_inner.stand_id 
					FROM kermesses_stands ks_inner
					JOIN kermesses k ON ks_inner.kermesse_id = k.id
					WHERE k.status = 'STARTED'
				)
			)
    `
	}
	err := s.db.Select(&stands, query)

	return stands, err
}

func (s *Repository) FindById(id int) (models.Stand, error) {
	stand := models.Stand{}
	query := "SELECT * FROM stands WHERE id=$1"
	err := s.db.Get(&stand, query, id)

	return stand, err
}

func (s *Repository) FindByUserId(userId int) (models.Stand, error) {
	stand := models.Stand{}
	query := "SELECT * FROM stands WHERE user_id=$1 LIMIT 1"
	err := s.db.Get(&stand, query, userId)

	return stand, err
}

func (s *Repository) Create(input map[string]interface{}) error {
	query := "INSERT INTO stands (user_id, name, description, type, price, stock) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := s.db.Exec(query, input["user_id"], input["name"], input["description"], input["type"], input["price"], input["stock"])

	return err
}

func (s *Repository) Update(id int, input map[string]interface{}) error {
	query := "UPDATE stands SET name=$1, description=$2, price=$3, stock=$4 WHERE id=$5"
	_, err := s.db.Exec(query, input["name"], input["description"], input["price"], input["stock"], id)

	return err
}

func (s *Repository) UpdateByUserId(userId int, input map[string]interface{}) error {
	query := "UPDATE stands SET name=$1, description=$2, price=$3, stock=$4 WHERE user_id=$5"
	_, err := s.db.Exec(query, input["name"], input["description"], input["price"], input["stock"], userId)

	return err
}

func (s *Repository) UpdateStock(id int, quantity int) error {
	query := "UPDATE stands SET stock=stock+$1 WHERE id=$2"
	_, err := s.db.Exec(query, quantity, id)

	return err
}
