package useCase

import (
	"database/sql"
)

// Старая структура для UserRegister
type PostgreUser struct {
	db *sql.DB
}

func NewPostgreUser(db *sql.DB) *PostgreUser {
	return &PostgreUser{db: db}
}

// Новая структура для UserLogin
type PostgreLogin struct {
	db *sql.DB
}

func NewPostgreLogin(db *sql.DB) *PostgreLogin {
	return &PostgreLogin{db: db}
}

// Новая структура для Refreshtoken
type PostgreRefresh struct {
	db *sql.DB
}

func NewPostgreRefresh(db *sql.DB) *PostgreRefresh {
	return &PostgreRefresh{db: db}
}
