package db

import (
	"database/sql"
	"database/sql/driver"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
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
	Id         int        `db:"id"`
	UserId     int        `db:"user_id"`
	Name       string     `db:"name"`
	PrivateKey PrivateKey `db:"private_key"`
	AllowedIps string     `db:"allowed_ips"`
	Num        int        `db:"num"` // used to generate device IP
}

type PrivateKey struct {
	wgtypes.Key
}

// Value implements Valuer interface
func (p *PrivateKey) Value() (driver.Value, error) {
	return driver.Value(p.String()), nil
}

// Scan implements Scanner interface
func (p *PrivateKey) Scan(src interface{}) error {
	key, err := wgtypes.ParseKey(src.(string))
	if err != nil {
		return err
	}

	*p = PrivateKey{key}

	return nil
}
