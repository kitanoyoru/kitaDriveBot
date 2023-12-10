package commands

import (
	"log"

	"github.com/kitanoyoru/kitaDriveBot/apps/service/internal/config"
	"github.com/spf13/cobra"
)

var cfg config.Config

var devFlag bool

var rootCmd = &cobra.Command{
	Use:   "effective-mobile-task",
	Short: "bla bla bla",
	Long:  "more bla bla bla",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if err := config.ReadConfig(&cfg); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(serverCommand)
}
