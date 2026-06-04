package main

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "bot",
		Short: "Telegram bot that stores PDFs in Google Drive",
	}

	rootCmd.AddCommand(newRunCommand())
	rootCmd.AddCommand(newAuthCommand())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("command failed")
	}
}
