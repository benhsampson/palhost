package services

import (
	"log"

	"github.com/stretchr/testify/assert"
)

func (s *DBTestSuite) SeedUsers(users []User) {
	store := UserStore{db: s.db}

	for _, user := range users {
		if _, err := store.CreateUser(&user); err != nil {
			log.Fatal(err)
		}
	}
}

const T_U string = "test"          // username
const T_P string = "test"          // password
const T_E string = "test@test.com" // email
const T_E_X string = "invalid"     // email (invalid)

func (s *DBTestSuite) TestCreateUser() {
	store := UserStore{db: s.db}

	id, err := store.CreateUser(&User{Username: T_U, Password: T_P, Email: T_E})
	assert.NoError(s.T(), err)

	user, err := store.GetUser("test")
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.ID, id)
	assert.Equal(s.T(), user.Username, T_U)
	assert.NotEqual(s.T(), user.Password, T_P, "password should be hashed")
	assert.Equal(s.T(), user.Email, T_E)
}

func (s *DBTestSuite) TestCreateUserUsernameCollision() {
	s.SeedUsers([]User{{Username: T_U, Password: T_P, Email: T_E}})

	store := UserStore{db: s.db}
	_, err := store.CreateUser(&User{Username: T_U, Password: T_P, Email: T_E})
	assert.Error(s.T(), err)
	assert.IsType(s.T(), ErrUsernameTaken{}, err)
}

func (s *DBTestSuite) TestCreateUserValidation() {
	store := UserStore{db: s.db}

	v := StructValidation[User]{callback: func(u *User) error { _, err := store.CreateUser(u); return err }, t: s.T()}

	v.Test(&User{Username: ""}, ValidationMap{"Username": {"required"}, "Email": {"required"}, "Password": {"required"}})

	v.Test(&User{Username: T_U, Email: T_E_X, Password: T_P}, ValidationMap{"Email": {"email"}})
}
