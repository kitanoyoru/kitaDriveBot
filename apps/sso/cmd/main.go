package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/kitanoyoru/kitaDriveBot/apps/sso/cmd/migrate"
)

var rootCmd = &cobra.Command{
	Short: "sso",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func main() {
	rootCmd.AddCommand(migrate.Command())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
