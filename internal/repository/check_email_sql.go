package repository

import (
	"context"
	"database/sql"
)

type Database struct {
	DB *sql.DB
}

// запрос в базу данных проверки существования email
func (db *Database) CheckEmailExists(email string) (bool, error) {
	var exists bool
	err := db.DB.QueryRowContext(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email).Scan(&exists)
	return exists, err
}
