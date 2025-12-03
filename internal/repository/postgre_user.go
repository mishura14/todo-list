package repository

import (
	"database/sql"
)

type PostgreUser struct {
	db *sql.DB
}

func NewPostgreUser(db *sql.DB) *PostgreUser {
	return &PostgreUser{db: db}
}
