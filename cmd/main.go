package main

import (
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cobra.Command{
		Short: "Simple wireguard web interface",
	}

	rootCmd.AddCommand(newStartCmd())
	rootCmd.AddCommand(newDbCmd())

	rootCmd.Execute()
}
