package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

const (
	subscribersTable             = "subscribers"
	servicesTable                = "services"
	subscribersServicesTable     = "subscribers_services"
	subscribersServiceNamesTable = "subscribers_service_names"
	messagesTable                = "messages"
)

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName)
	db, err := sqlx.Open("postgres", url)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
