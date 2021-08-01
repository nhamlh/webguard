package db

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB

func init () {
	DB = initDb()
}

func initDb() *sqlx.DB {
	dbPath := "/tmp/wg-dash.db"

	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	return db
}

func MigrateSchema() {
	DB.MustExec(schema)

	// Seeding
	DB.MustExec(`
INSERT INTO users(email,password,is_admin,auth_type)
values("nham", "$2a$14$VxNB2aRwQj0eueo.1g25YOtnga/9AmSxAeHX5hnXpdDszat1COob2", 1, 0)
ON CONFLICT DO NOTHING;

INSERT INTO devices(user_id,name,private_key,allowed_ips)
values(1, "one", "EK0si+W7Pu2UoMYUE83/Dj6XNVrgROVACVpBlVOwEV0=", "0.0.0.0/0");
INSERT INTO devices(user_id,name,private_key,allowed_ips)
values(1, "two", "gKM16ZjyrnAY5OjsEwLqX7FXu6DySZwmrEsjOFq2DW8=", "1.1.1.1/1");

INSERT INTO users(email,password,is_admin,auth_type)
values("test", "$2a$14$VxNB2aRwQj0eueo.1g25YOtnga/9AmSxAeHX5hnXpdDszat1COob2", 1, 0)
ON CONFLICT DO NOTHING;

INSERT INTO devices(user_id,name,private_key,allowed_ips)
values(2, "four", "8O9JcFJasK2zUvtC77sD2nWzVnIl3lOibIcKG3Z3238=", "0.0.0.0/0");
`)
}
