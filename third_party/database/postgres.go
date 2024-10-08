package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func NewPostgres(config PostgresConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		config.Host, config.Port, config.User, config.Password, config.Name,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
