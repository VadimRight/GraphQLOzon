package repository

import "database/sql"

type UserRepository struct {
	Conn *sql.DB
	Collection string
}
