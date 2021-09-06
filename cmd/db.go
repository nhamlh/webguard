package main

import (
	"bytes"
	"html/template"
	"log"

	"github.com/nhamlh/wg-dash/pkg/db"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

func newDbCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "db",
		Short: "Database related operations",
	}

	cmd.AddCommand(newDBMigrateCmd())
	cmd.AddCommand(newDBSeedCmd())

	return cmd
}

func newDBMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate database schema",
		Run: func(cmd *cobra.Command, args []string) {
			if err := goose.SetDialect("sqlite3"); err != nil {
				panic(err)
			}

			goose.SetBaseFS(db.Migrations)

			if err := goose.Up(db.DB.DB, "migrations"); err != nil {
				panic(err)
			}
		},
	}

	return cmd
}

func newDBSeedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Database data seeding",
		Run: func(cmd *cobra.Command, args []string) {
			tmpl, _ := template.New("sql").Parse(`
insert into users(email,password,is_admin,auth_type)
values("nham", "{{ $.password }}", 1, 0)
on conflict do nothing;

insert into devices(user_id,name,private_key,allowed_ips)
values(1, "one", "ek0si+w7pu2uomyue83/dj6xnvrgrovacvpblvowev0=", "0.0.0.0/0");
insert into devices(user_id,name,private_key,allowed_ips)
values(1, "two", "gkm16zjyrnay5ojsewlqx7fxu6dyszwmresjofq2dw8=", "1.1.1.1/1");

insert into users(email,password,is_admin,auth_type)
values("test", "{{ $.password }}", 1, 0)
on conflict do nothing;

insert into devices(user_id,name,private_key,allowed_ips)
values(2, "four", "8o9jcfjask2zuvtc77sd2nwzvnil3loibickg3z3238=", "0.0.0.0/0");
`)
			password, err := bcrypt.GenerateFromPassword([]byte("abc"), bcrypt.DefaultCost)
			if err != nil {
				log.Fatal(err)
			}

			sql := bytes.NewBufferString("")
			tmpl.Execute(sql, map[string]interface{}{"password": string(password)})

			db.DB.MustExec(sql.String())
		},
	}

	return cmd
}
