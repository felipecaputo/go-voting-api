package user

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db sqlx.DB
}

const insertUserSQL = "INSERT INTO user (id, name, email, password, is_admin) VALUES (:id,:name,:email,:password,:is_admin)"
const getUserDataSQL = "SELECT id, name, email, is_admin FROM user WHERE id = ?"
const updateUserDataSQL = "UPDATE user SET %s WHERE id =:id"
const deleteUserSQL = "DELETE FROM user where id = ?"

func hashAndSaltPwd(password string) (string, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(pwd), err
}

func (u *UserRepository) updateUser(user User, fields []string) error {
	var updateFields = make([]string, len(fields))
	for i, f := range fields {
		updateFields[i] = fmt.Sprintf("%s=:%s", f, f)
	}

	qry := fmt.Sprintf(updateUserDataSQL, strings.Join(updateFields, ", "))
	_, err := u.db.NamedExec(qry, user)

	return err
}

// CreateUser insert the giver user into database with salted password
// and return the the inserted user with ID
func (u *UserRepository) CreateUser(user User) (*User, error) {
	user.ID = uuid.New().String()

	var err error

	user.Password, err = hashAndSaltPwd(user.Password)

	if err != nil {
		return nil, err
	}

	if _, err = u.db.NamedExec(insertUserSQL, user); err != nil {
		return nil, errors.New("error while creating user.")
	}

	user.Password = ""
	return &user, nil
}

func (u *UserRepository) Get(id string) (User, error) {
	user := User{}
	err := u.db.Get(&user, getUserDataSQL, id)
	return user, err
}

func (u *UserRepository) Update(user User) error {
	return u.updateUser(user, []string{"name", "email"})
}

func (u *UserRepository) Delete(id string) error {
	_, err := u.db.Exec(deleteUserSQL, id)
	return err
}
