package main

import (
	"fmt"
	"github.com/nhamlh/wg-dash/pkg/config"
	"github.com/nhamlh/wg-dash/pkg/web"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"time"
)

func newStartCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start wg-dash server",
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

			router := web.NewRouter()

			srv := &http.Server{
				Handler: router,
				Addr:    fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Web.ListenPort),
				// Good practice: enforce timeouts for servers you create!
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}

			fmt.Println(fmt.Sprintf("Web server is listening at %s:%d", cfg.Hostname, cfg.Web.ListenPort))
			srv.ListenAndServe()
		},
	}

	return cmd
}
