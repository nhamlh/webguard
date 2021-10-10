package db

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type User struct {
	Id        int            `db:"id"`
	Email     string         `db:"email"`
	Password  sql.NullString `db:"password"`
	AuthType  AuthType       `db:"auth_type"`
	IsAmdin   bool           `db:"is_admin"`
	LastLogin Time           `db:"last_login"`
}

func AllUsers(db sqlx.DB) (users []User) {
	db.Select(&users, "SELECT * FROM users")

	return
}

func GetUserById(id int, db sqlx.DB) (User, bool) {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	if user == (User{}) || err != nil {
		return User{}, false
	}

	return user, true
}

func GetUserByEmail(email string, db sqlx.DB) (User, bool) {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE email=$1", email)
	if user == (User{}) || err != nil {
		return User{}, false
	}

	return user, true
}

// Save persists user to database
func (u *User) Save(db sqlx.DB) error {
	var tmp User
	db.Get(&tmp, `SELECT * FROM users WHERE id=$1`, u.Id)

	if tmp == (User{}) {
		log.Println("Inserting new user")
		_, err := db.Exec(`
INSERT INTO
users(email, password, auth_type, is_admin, last_login)
values ($1,$2,$3,$4,$5)
`,
			u.Email,
			u.Password,
			u.AuthType,
			u.IsAmdin,
			u.LastLogin)

		return err
	} else {
		log.Println("Updating existing user")
		_, err := db.Exec(`
UPDATE users SET email=$1, password=$2, is_admin=$3, last_login=$4
WHERE id=$5`,
			u.Email,
			u.Password,
			u.IsAmdin,
			u.LastLogin,
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

func (u *User) RecordLogin() {
	u.LastLogin = Time{time.Now()}
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

type Time struct {
	time.Time
}

// Value implements Valuer interface
func (t Time) Value() (driver.Value, error) {
	return driver.Value(t.Format(time.RFC3339)), nil
}

// Scan implements Scanner interface
func (t *Time) Scan(src interface{}) error {
	if src == nil {
		*t = Time{time.Time{}}
		return nil
	}

	strVal, ok := src.(string)
	if !ok {
		return errors.New("Field is not a type of string")
	}

	parsedTime, err := time.Parse(time.RFC3339, strVal)
	if err != nil {
		return err
	}

	*t = Time{parsedTime}

	return nil
}
