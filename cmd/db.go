package main

import (
	"github.com/spf13/cobra"
	"github.com/nhamlh/wg-dash/pkg/db"
)

func newDbCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use: "db",
		Short: "Database related operations",
	}

	cmd.AddCommand(newDBMigrateCmd())

	return cmd
}

func newDBMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "migrate",
		Short: "Migrate database schema",
		Run: func(cmd *cobra.Command, args []string) {
			db.MigrateSchema()
		},
	}

	return cmd
}
