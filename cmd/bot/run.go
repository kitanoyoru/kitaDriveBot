package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/kitanoyoru/kitaDriveBot/internal/app"
	"github.com/spf13/cobra"
)

func newRunCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Start the Telegram bot",
		RunE: func(_ *cobra.Command, _ []string) error {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer stop()
			return app.RunBot(ctx)
		},
	}
}
