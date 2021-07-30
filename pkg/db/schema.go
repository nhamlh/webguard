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
)
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

// func (p *Profile) AddDevice(d Device) error {
// 	append(d, p.Devices)
// 	return nil
// }

// func (p *Profile) DeleteDevice(id string) error {
// 	for d := range p.Devices {
// 		if d.Id == id {
// 			delete(d, p.Devices)
// 		}
// 	}
// 	return nil
// }

const (
	LinuxDevice = iota
	MacDevice
	WindowsDevice
	AndroidDevice
	IosDevice
)

type DeviceType int

type Device struct {
	Id         int
	Name       string
	PrivateKey string
	Type       DeviceType
	AllowedIps string
}
