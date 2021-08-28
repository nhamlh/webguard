package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	// "github.com/nhamlh/wg-dash/pkg/config"
)

var DB *sqlx.DB

func init() {
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

// func migrateschema() {
// 	// seeding
// 	db.mustexec(`
// insert into users(email,password,is_admin,auth_type)
// values("nham", "$2a$14$vxnb2arwqj0eueo.1g25yotnga/9amsxaehx5hnxpddszat1coob2", 1, 0)
// on conflict do nothing;

// insert into devices(user_id,name,private_key,allowed_ips)
// values(1, "one", "ek0si+w7pu2uomyue83/dj6xnvrgrovacvpblvowev0=", "0.0.0.0/0");
// insert into devices(user_id,name,private_key,allowed_ips)
// values(1, "two", "gkm16zjyrnay5ojsewlqx7fxu6dyszwmresjofq2dw8=", "1.1.1.1/1");

// insert into users(email,password,is_admin,auth_type)
// values("test", "$2a$14$vxnb2arwqj0eueo.1g25yotnga/9amsxaehx5hnxpddszat1coob2", 1, 0)
// on conflict do nothing;

// insert into devices(user_id,name,private_key,allowed_ips)
// values(2, "four", "8o9jcfjask2zuvtc77sd2nwzvnil3loibickg3z3238=", "0.0.0.0/0");
// `)
// }
