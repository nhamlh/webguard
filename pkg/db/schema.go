package db

import (
	"database/sql"
)

const (
	StaticAuth = iota
	SSOAuth
)

type User struct {
	Id       int            `db:"id"`
	Email    string         `db:"email"`
	Password sql.NullString `db:"password"`
	AuthType int            `db:"auth_type"`
	IsAmdin  bool           `db:"is_admin"`
}

func (u *User) All()     {}
func (u *User) ById()    {}
func (u *User) ByEmail() {}
func (u *User) Save()    {}

