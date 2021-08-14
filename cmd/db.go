package main

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/nhamlh/wg-dash/pkg/db"
	"github.com/spf13/cobra"
	"log"
)

func newDbCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "db",
		Short: "Database related operations",
	}

	cmd.AddCommand(newDBMigrateCmd())

	return cmd
}

func newDBMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate database schema",
		Run: func(cmd *cobra.Command, args []string) {
			driver, err := sqlite3.WithInstance(db.DB.DB, &sqlite3.Config{})
			if err != nil {
				log.Fatal(err)
			}

			m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "sqlite3", driver)
			if err != nil {
				log.Fatal(err)
			}

			m.Steps(2)
		},
	}

	return cmd
}
