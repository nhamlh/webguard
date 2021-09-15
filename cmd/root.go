package main

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cobra.Command{
		Short: "Simple wireguard web interface",
	}

	rootCmd.PersistentFlags().StringP("config", "c", "", "Path to config file")

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(dbCmd)
	rootCmd.AddCommand(genConfigCmd)

	rootCmd.Execute()
}
