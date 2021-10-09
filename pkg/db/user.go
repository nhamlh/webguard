package db

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type User struct {
	Id       int            `db:"id"`
	Email    string         `db:"email"`
	Password sql.NullString `db:"password"`
	AuthType AuthType       `db:"auth_type"`
	IsAmdin  bool           `db:"is_admin"`
}

func AllUsers(db sqlx.DB) (users []User) {
	db.Select(&users, "SELECT * FROM users")

	return
}

func (u *User) ById()    {}
func (u *User) ByEmail() {}

// Save persists user to database
func (u *User) Save(db sqlx.DB) error {
	var tmp User
	db.Get(&tmp, `SELECT * FROM users WHERE id=$1`, u.Id)

	if tmp == (User{}) {
		log.Println("Inserting new user")
		_, err := db.Exec(`
INSERT INTO
users(email, password, auth_type, is_admin)
values ($1,$2,$3,$4)
`,
			u.Email,
			u.Password,
			u.AuthType,
			u.IsAmdin)

		return err
	} else {
		log.Println("Updating existing user")
		_, err := db.Exec(`
UPDATE users SET email=$1, password=$2, is_admin=$3 WHERE id=$4`,
			u.Email,
			u.Password,
			u.IsAmdin,
			u.Id)

		return err
	}
}

// PasswdMatched return true if provided passwd match
// user's hashed passwd
func (u *User) PasswdMatched(passwd []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password.String), passwd)
	if err != nil {
		return false
	}

	return true
}

func (u *User) NewPasswd(new string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(new), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Cannot hash new password: %v", err)
	}

	u.Password = sql.NullString{String: string(hashedPass), Valid: true}

	return nil
}

type AuthType int64

const (
	StaticAuth AuthType = iota
	SSOAuth
)

// Value implements Valuer interface
func (a *AuthType) Value() (driver.Value, error) {
	return driver.Value(*a), nil
}

// Scan implements Scanner interface
func (a *AuthType) Scan(src interface{}) error {
	if src.(int64) <= 1 {
		*a = AuthType(src.(int64))
		return nil
	} else {
		return errors.New("Invalid authentication type")
	}
}
