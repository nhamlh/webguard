package db

import (
	"embed"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var Migrations embed.FS

var DB *sqlx.DB

func init() {
	DB = initDb()
}

func initDb() *sqlx.DB {
	dbPath := "./webguard.db"

	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}

	return db
}
