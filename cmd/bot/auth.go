package main

import (
	"github.com/kitanoyoru/kitaDriveBot/internal/app"
	"github.com/spf13/cobra"
)

func newAuthCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "auth",
		Short: "Authorize Google Drive access (one-time setup)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.RunAuth(cmd.Context())
		},
	}
}
