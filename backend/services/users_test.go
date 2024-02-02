package services

import (
	"log"

	"github.com/stretchr/testify/assert"
)

func (s *DBTestSuite) SeedUsers(users []User) {
	store := UsersService{db: s.db}

	for _, user := range users {
		if _, err := store.CreateUser(&user); err != nil {
			log.Fatal(err)
		}
	}
}

const T_U string = "test"             // username
const T_U_2 string = "test2"          // username 2
const T_P string = "abcdefgh"         // password
const T_P_2 string = "hijklmnop"      // password
const T_P_X string = "short"          // password (invalid)
const T_E string = "test@test.com"    // email
const T_E_2 string = "test2@test.com" // email
const T_E_X string = "invalid"        // email (invalid)

func (s *DBTestSuite) TestCreateUser() {
	store := UsersService{db: s.db}

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

	store := UsersService{db: s.db}
	_, err := store.CreateUser(&User{Username: T_U, Password: T_P, Email: T_E})
	assert.Error(s.T(), err)
	assert.IsType(s.T(), ErrUsernameTaken{}, err)
}

func (s *DBTestSuite) TestCreateUserValidation() {
	store := UsersService{db: s.db}

	v := StructValidation[User]{callback: func(u *User) error { _, err := store.CreateUser(u); return err }, t: s.T()}

	v.Test(&User{Username: ""}, ValidationMap{"Username": {"required"}, "Email": {"required"}, "Password": {"required"}})
	v.Test(&User{Username: T_U, Email: T_E_X, Password: T_P}, ValidationMap{"Email": {"email"}})
}

func (s *DBTestSuite) TestSignIn() {
	s.SeedUsers([]User{{Username: T_U, Password: T_P, Email: T_E}})

	store := UsersService{db: s.db}
	user, err := store.SignIn(T_U, T_P)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.Username, T_U)
	assert.NotEqual(s.T(), user.Password, T_P, "password should be hashed")
}

func (s *DBTestSuite) TestInvalidSignIn() {
	s.SeedUsers([]User{{Username: T_U, Password: T_P, Email: T_E}})

	store := UsersService{db: s.db}
	user, err := store.SignIn(T_U, T_P_X)
	assert.Error(s.T(), err)
	assert.IsType(s.T(), ErrInvalidPassword{}, err)
	assert.Nil(s.T(), user)
}

func (s *DBTestSuite) TestUpdateUser() {
	s.SeedUsers([]User{{Username: T_U, Password: T_P, Email: T_E}})

	store := UsersService{db: s.db}
	user, err := store.UpdateUser(T_U, &User{Username: T_U_2, Email: T_E_2})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.Username, T_U_2)
	assert.Equal(s.T(), user.Email, T_E_2)

	user, err = store.GetUser(T_U_2)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.Username, T_U_2)
	assert.Equal(s.T(), user.Email, T_E_2)
}

func (s *DBTestSuite) TestUpdateUserValidation() {
	store := UsersService{db: s.db}

	v := StructValidation[User]{callback: func(u *User) error { _, err := store.UpdateUser(T_U, u); return err }, t: s.T()}

	v.Test(&User{Email: ""}, ValidationMap{"Username": {"required"}, "Email": {"required"}})
	v.Test(&User{Username: T_U, Email: T_E_X}, ValidationMap{"Email": {"email"}})
}

func (s *DBTestSuite) TestUpdateUserNotFound() {
	store := UsersService{db: s.db}

	user, err := store.UpdateUser(T_U, &User{Username: T_U_2, Email: T_E_2})
	assert.Error(s.T(), err)
	assert.IsType(s.T(), ErrUserNotFound{}, err)
	assert.Nil(s.T(), user)
}

func (s *DBTestSuite) TestUpdatePassword() {
	s.SeedUsers([]User{{Username: T_U, Password: T_P, Email: T_E}})

	store := UsersService{db: s.db}
	err := store.UpdatePassword(T_U, &UserChangePassword{CurrentPassword: T_P, NewPassword: T_P_2, ConfirmPassword: T_P_2})
	assert.NoError(s.T(), err)

	_, err = store.SignIn(T_U, T_P)
	assert.Error(s.T(), err)

	_, err = store.SignIn(T_U, T_P_2)
	assert.NoError(s.T(), err)
}

func (s *DBTestSuite) TestUpdatePasswordValidation() {
	store := UsersService{db: s.db}

	v := StructValidation[UserChangePassword]{callback: func(p *UserChangePassword) error { return store.UpdatePassword(T_U, p) }, t: s.T()}

	v.Test(&UserChangePassword{CurrentPassword: ""}, ValidationMap{"CurrentPassword": {"required"}, "NewPassword": {"required"}, "ConfirmPassword": {"required"}})
	v.Test(&UserChangePassword{CurrentPassword: T_P, NewPassword: T_P_X, ConfirmPassword: T_P_X}, ValidationMap{"NewPassword": {"min"}})
	v.Test(&UserChangePassword{CurrentPassword: T_P, NewPassword: T_P, ConfirmPassword: T_P_2}, ValidationMap{"ConfirmPassword": {"eqfield"}})
}
