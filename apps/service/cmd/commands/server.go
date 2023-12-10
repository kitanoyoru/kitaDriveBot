package commands

import (
	"log"

	"github.com/kitanoyoru/kitaDriveBot/apps/service/internal/app"
	"github.com/spf13/cobra"
)

var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long:  "Used to launch and run the service with your env variables as a config.",
	Run: func(cmd *cobra.Command, args []string) {
		app, err := app.NewApp(&cfg)
		if err != nil {
			log.Fatal(err)
		}

		if err := app.Run(); err != nil {
			log.Fatal(err)
		}
	},
}
