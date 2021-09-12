package main

import (
	"log"

	"github.com/nhamlh/webguard/pkg/config"
	models "github.com/nhamlh/webguard/pkg/db"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
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
			cfgFile, err := cmd.Flags().GetString("config")
			if err != nil {
				log.Fatal(err)
			}

			var cfg *config.Config
			if cfgFile == "" {
				cfg = &config.DefaultConfig
			} else {
				cfg = config.Load(cfgFile)
			}

			db := models.InitDb(cfg.DbPath)

			if err := goose.SetDialect("sqlite3"); err != nil {
				panic(err)
			}

			goose.SetBaseFS(models.Migrations)

			if err := goose.Up(db.DB, "migrations"); err != nil {
				panic(err)
			}
		},
	}

	return cmd
}
