package main

import (
	"encoding/json"
	"fmt"

	"github.com/nhamlh/webguard/pkg/config"
	"github.com/spf13/cobra"
)

var genConfigCmd = &cobra.Command{
	Use:   "genconf",
	Short: "Generate config file and print to stdout",
	RunE: func(cmd *cobra.Command, args []string) error {
		ret, err := json.MarshalIndent(config.DefaultConfig, "", "  ")
		if err != nil {
			return fmt.Errorf("Cannot marshal config struct: %v", err)
		}

		fmt.Println(string(ret))
		return nil
	},
}
