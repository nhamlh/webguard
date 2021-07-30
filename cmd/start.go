package main

import (
	// "github.com/nhamlh/wg-dash/pkg/config"
	"github.com/nhamlh/wg-dash/pkg/web"
	"github.com/spf13/cobra"
)
func newStartCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use: "start",
		Short: "Start wg-dash server",
		Run: func(cmd *cobra.Command, args []string) {
			// cfg := config.Load()
			web.StartServer()
		},
	}

	return cmd
}

