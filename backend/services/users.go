package services

import (
	"database/sql"

	"github.com/go-playground/validator/v10"
)

type UserStore struct {
	db *sql.DB
}

type User struct {
	ID       int    `validate:"required"`
	Username string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type ErrUsernameTaken struct {
	error
}

var val = validator.New()

func (s *UserStore) CreateUser(u *User) (int, error) {
	if err := val.StructExcept(u, "ID"); err != nil {
		valErrors := err.(validator.ValidationErrors)
		return 0, valErrors
	}
	query := `
	SELECT id FROM USERS WHERE username = $1;`
	row := s.db.QueryRow(query, u.Username)
	var userId int
	if err := row.Scan(&userId); err != nil {
		if err != sql.ErrNoRows {
			return 0, err
		}
	} else {
		return 0, ErrUsernameTaken{err}
	}
	passwordHash, err := HashPassword(u.Password)
	if err != nil {
		return 0, err
	}
	query = `
	INSERT INTO users (username, password, email)
	VALUES ($1, $2, $3)
	RETURNING ID;`
	row = s.db.QueryRow(query, u.Username, passwordHash, u.Email)
	if err := row.Scan(&userId); err != nil {
		return 0, err
	}
	return userId, nil
}

func (s *UserStore) GetUser(username string) (*User, error) {
	query := `
	SELECT id, username, password, email FROM users WHERE username = $1;`
	row := s.db.QueryRow(query, username)
	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return &user, nil
}
