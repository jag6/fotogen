package models

import "database/sql"

type User struct {
	ID           int
	Email        string
	Username     string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}
