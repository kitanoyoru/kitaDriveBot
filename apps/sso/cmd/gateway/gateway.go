package gateway

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/app/gateway"
	"github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/config"
)

func Command() *cobra.Command {
	return &cobra.Command{
		Use: "gateway",
		RunE: func(cmd *cobra.Command, args []string) error {
			app, err := gateway.NewApp(config.Config{
				GrpcConfig: config.GrpcConfig{
					GRPCEndpoint: os.Getenv("GRPC_ENDPOINT"),
				},
				DatabaseConfig: config.DatabaseConfig{
					ConnectionString: os.Getenv("DATABASE_CONNECTION_STRING"),
				},
				LoggerConfig: config.LoggerConfig{
					LogLevel: os.Getenv("LOG_LEVEL"),
				},
			})
			if err != nil {
				return errors.Wrap(err, "failed to create app")
			}

			go func() {
				if err := app.Run(); err != nil {
					log.Err(err).Send()
				}
			}()

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)
			<-c

			fmt.Println("Shutting down server...")
			err = app.Close()
			if err != nil {
				log.Fatal().Err(err).Send()
			}

			return nil
		},
	}
}
