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

type Device struct {
	Id         int    `db:"id"`
	UserId     int    `db:"user_id"`
	Name       string `db:"name"`
	PrivateKey string `db:"private_key"`
	AllowedIps string `db:"allowed_ips"`
}
