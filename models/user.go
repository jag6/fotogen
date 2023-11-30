package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Email        string
	Username     string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(email, username, password string) (*User, error) {
	email = strings.ToLower(email)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	passwordHash := string(hashedBytes)

	user := User{
		Email:        email,
		Username:     username,
		PasswordHash: passwordHash,
	}
	row := us.DB.QueryRow(`
		INSERT INTO users (email, username, password_hash)
		VALUES ($1, $2, $3) RETURNING id`, email, username, passwordHash)
	err = row.Scan(&user.ID) //user.Username
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			var columnName = pgError.ColumnName
			if pgError.Code == pgerrcode.UniqueViolation && columnName == "email" {
				return nil, ErrEmailTaken
			} else if pgError.Code == pgerrcode.UniqueViolation && columnName == "username" {
				return nil, ErrUsernameTaken
			}
		}
		return nil, fmt.Errorf("create user :%w", err)
	}
	return &user, nil
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)
	user := User{
		Email: email,
	}
	row := us.DB.QueryRow(`
		SELECT id, password_hash 
		FROM users WHERE email=$1`, email)
	err := row.Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
	return &user, nil
}

func (us *UserService) UpdatePassword(userID int, password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	passwordHash := string(hashedBytes)
	_, err = us.DB.Exec(`
		UPDATE users
		SET password_hash = $2
		WHERE id = $1;`, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}

func (us *UserService) SearchByUsername(q string) ([]User, error) {
	rows, err := us.DB.Query(`
		SELECT id, username 
		FROM users
		WHERE username = $1;`, q)
	if err != nil {
		return nil, fmt.Errorf("query users by username: %w", err)
	}
	var users []User
	for rows.Next() {
		user := User{
			Username: q,
		}
		err = rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return nil, fmt.Errorf("query users by username: %w", err)
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("query users by username: %w", err)
	}
	return users, nil
}
