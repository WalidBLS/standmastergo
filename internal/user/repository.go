package user

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"standmaster/internal/models"
)

type UserRepository interface {
	FindAll(filters map[string]interface{}) ([]models.UserBasic, error)
	FindAllChildren(id int, filters map[string]interface{}) ([]models.UserBasic, error)
	FindById(id int) (models.User, error)
	FindByEmail(email string) (models.User, error)
	Create(input map[string]interface{}) error
	Update(id int, input map[string]interface{}) error
	UpdateCredit(id int, amount int) error
	HasStand(id int) (bool, error)
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (s *Repository) FindAll(filters map[string]interface{}) ([]models.UserBasic, error) {
	users := []models.UserBasic{}
	query := `
		SELECT DISTINCT
			u.id AS id,
			u.name AS name,
			u.email AS email,
			u.role AS role,
			u.credit AS credit
		FROM users u
		FULL OUTER JOIN kermesses_users ku ON u.id = ku.user_id
		WHERE 1=1
	`
	if filters["kermesse_id"] != nil {
		query += fmt.Sprintf(" AND ku.kermesse_id = %v", filters["kermesse_id"])
	}
	err := s.db.Select(&users, query)

	return users, err
}

func (s *Repository) FindAllChildren(id int, filters map[string]interface{}) ([]models.UserBasic, error) {
	users := []models.UserBasic{}
	query := `
		SELECT DISTINCT
			u.id AS id,
			u.name AS name,
			u.email AS email,
			u.role AS role,
			u.credit AS credit
		FROM users u
		FULL OUTER JOIN kermesses_users ku ON u.id = ku.user_id
		WHERE u.role=$1 AND u.parent_id=$2
	`
	if filters["kermesse_id"] != nil {
		query += fmt.Sprintf(" AND ku.kermesse_id = %v", filters["kermesse_id"])
	}
	err := s.db.Select(&users, query, models.UserRoleChild, id)

	return users, err
}

func (s *Repository) FindById(id int) (models.User, error) {
	user := models.User{}
	query := "SELECT * FROM users WHERE id=$1"
	err := s.db.Get(&user, query, id)

	return user, err
}

func (s *Repository) FindByEmail(email string) (models.User, error) {
	user := models.User{}
	query := "SELECT * FROM users WHERE email=$1"
	err := s.db.Get(&user, query, email)

	return user, err
}

func (s *Repository) Create(input map[string]interface{}) error {
	query := "INSERT INTO users (parent_id, name, email, password, role) VALUES ($1, $2, $3, $4, $5)"
	_, err := s.db.Exec(query, input["parent_id"], input["name"], input["email"], input["password"], input["role"])

	return err
}

func (s *Repository) Update(id int, input map[string]interface{}) error {
	query := "UPDATE users SET password=$1 WHERE id=$2"
	_, err := s.db.Exec(query, input["new_password"], id)

	return err
}

func (s *Repository) UpdateCredit(id int, amount int) error {
	query := "UPDATE users SET credit=credit+$1 WHERE id=$2"
	_, err := s.db.Exec(query, amount, id)

	return err
}

func (s *Repository) HasStand(id int) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM stands WHERE user_id=$1"
	err := s.db.Get(&count, query, id)

	return count >= 1, err
}
