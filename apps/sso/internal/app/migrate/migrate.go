package migrate

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb" // db driver
	_ "github.com/golang-migrate/migrate/v4/database/postgres"    // db driver
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"

	"github.com/kitanoyoru/kitaDriveBot/apps/sso/pkg/migrations"
	"github.com/kitanoyoru/kitaDriveBot/libs/app"
)

func NewApp(sqlConnectionString string) (app.App, error) {
	s := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})

	d, err := bindata.WithInstance(s)
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithSourceInstance("go-bindata", d, sqlConnectionString)
	if err != nil {
		return nil, err
	}

	return &migrateApp{
		m: m,
	}, nil
}

type migrateApp struct {
	m *migrate.Migrate
}

func (app *migrateApp) Run() error {
	err := app.m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}

		return err
	}

	return nil
}

func (app *migrateApp) Close() error {
	return nil
}
