package db

import (
	"database/sql"
)

var schema = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT NOT NULL UNIQUE,
	password TEXT,
	is_admin TEXT NOT NULL,
	auth_type TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS devices (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	name TEXT NOT NULL UNIQUE,
	private_key TEXT NOT NULL,
	allowed_ips TEXT NOT NULL,

	FOREIGN KEY(user_id) REFERENCES users(id)
);
`

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
