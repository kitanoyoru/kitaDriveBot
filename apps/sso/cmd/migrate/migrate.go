package up

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/app/migrate"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use: "up",
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := migrate.NewApp(os.Getenv("SQL_CONNECTION_STRING"))
			if err != nil {
				return errors.Wrap(err, "failed to create app")
			}

			return app.Run()
		},
	}
}
