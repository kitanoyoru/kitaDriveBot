package migrator

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/kitanoyoru/kitaDriveBot/apps/sso/internal/app/migrate"
)

func Migrator(sqlConnectionString string) error {
	migrateApp, err := migrate.NewApp(strings.Replace(sqlConnectionString, "postgres", "cockroachdb", 1))
	if err != nil {
		return errors.Wrap(err, "failed to migrate up")
	}

	return migrateApp.Run()
}
