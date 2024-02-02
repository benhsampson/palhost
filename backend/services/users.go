package services

import (
	"database/sql"
	"errors"

	"github.com/go-playground/validator/v10"
)

type UsersService struct {
	db *sql.DB
}

type User struct {
	ID       int    `validate:"required"`
	Username string `validate:"required"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

type UserChangePassword struct {
	CurrentPassword string `validate:"required"`
	NewPassword     string `validate:"required,min=8"`
	ConfirmPassword string `validate:"required,eqfield=NewPassword"`
}

type ErrUserNotFound struct {
	error
}

type ErrUsernameTaken struct {
	error
}

var val = validator.New()

func (s *UsersService) GetUser(username string) (*User, error) {
	query := `
	SELECT id, username, password, email FROM users WHERE username = $1;`
	row := s.db.QueryRow(query, username)
	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound{err}
		}
		return nil, err
	}
	return &user, nil
}

func (s *UsersService) CreateUser(u *User) (int, error) {
	if err := val.StructPartial(u, "Username", "Email", "Password"); err != nil {
		valErrors := err.(validator.ValidationErrors)
		return 0, valErrors
	}
	_, err := s.GetUser(u.Username)
	if err != nil {
		if !errors.As(err, &ErrUserNotFound{}) {
			return 0, err
		}
	} else {
		return 0, ErrUsernameTaken{}
	}
	passwordHash, err := HashPassword(u.Password)
	if err != nil {
		return 0, err
	}
	query := `
	INSERT INTO users (username, password, email)
	VALUES ($1, $2, $3)
	RETURNING ID;`
	row := s.db.QueryRow(query, u.Username, passwordHash, u.Email)
	var userId int
	if err := row.Scan(&userId); err != nil {
		return 0, err
	}
	return userId, nil
}

type ErrInvalidPassword struct {
	error
}

func (s *UsersService) SignIn(username, password string) (*User, error) {
	user, err := s.GetUser(username)
	if err != nil {
		return nil, err
	}
	if !CheckPasswordHash(password, user.Password) {
		return nil, ErrInvalidPassword{}
	}
	return user, nil
}

func (s *UsersService) UpdateUser(username string, u *User) (*User, error) {
	if err := val.StructPartial(u, "Username", "Email"); err != nil {
		valErrors := err.(validator.ValidationErrors)
		return nil, valErrors
	}
	if _, err := s.GetUser(username); err != nil {
		return nil, err
	}
	query := `
	UPDATE users 
	SET username = $2, email = $3
	WHERE username = $1
	RETURNING id, username, email;`
	row := s.db.QueryRow(query, username, u.Username, u.Email)
	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Email); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UsersService) UpdatePassword(username string, p *UserChangePassword) error {
	if err := val.Struct(p); err != nil {
		return err.(validator.ValidationErrors)
	}
	user, err := s.GetUser(username)
	if err != nil {
		return err
	}
	if !CheckPasswordHash(p.CurrentPassword, user.Password) {
		return ErrInvalidPassword{}
	}
	passwordHash, err := HashPassword(p.NewPassword)
	if err != nil {
		return err
	}
	query := `
	UPDATE users 
	SET password = $2
	WHERE username = $1;`
	if _, err = s.db.Exec(query, username, passwordHash); err != nil {
		return err
	}
	return nil
}
