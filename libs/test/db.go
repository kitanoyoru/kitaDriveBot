package test

import (
	"context"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/kitanoyoru/kitaDriveBot/libs/test/testcontainers"
)

type DBTest struct {
	dbPool *pgxpool.Pool
}

func (dbTest *DBTest) CreateDB(t *testing.T, props ...MigrationProperty) (*pgxpool.Pool, error) {
	if dbTest.dbPool == nil {
		sqlConnectionString, err := getSQLConnectionString(t)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create cockroach container")
		}

		dbPool, err := getDBPool(sqlConnectionString)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create db pool")
		}

		dbTest.dbPool = dbPool
	}

	migrationsProps := &MigrationProperties{}
	for _, p := range props {
		p(migrationsProps)
	}

	err := optimizeDB(dbTest.dbPool)
	if err != nil {
		return nil, errors.Wrap(err, "failed to optimize db")
	}

	return createDB(dbTest.dbPool, migrationsProps)
}

func optimizeDB(dbPool *pgxpool.Pool) error {
	queries := []string{
		"SET CLUSTER SETTING kv.range_merge.queue_interval = '50ms';",
		"SET CLUSTER SETTING jobs.registry.interval.gc = '30s';",
		"SET CLUSTER SETTING jobs.registry.interval.cancel = '180s';",
		"SET CLUSTER SETTING jobs.retention_time = '15s';",
		"SET CLUSTER SETTING sql.stats.automatic_collection.enabled = false;",
		"SET CLUSTER SETTING kv.range_split.by_load_merge_delay = '5s';",
		"ALTER RANGE default CONFIGURE ZONE USING \"gc.ttlseconds\" = 600;",
		"ALTER DATABASE system CONFIGURE ZONE USING \"gc.ttlseconds\" = 600;",
	}

	for _, query := range queries {
		_, err := dbPool.Exec(context.Background(), query)
		if err != nil {
			return errors.Wrap(err, "failed to execute query: "+query)
		}
	}

	return nil
}

func getSQLConnectionString(t *testing.T) (string, error) {
	if v, ok := os.LookupEnv("APP_INTEGRATION_TESTS_SQL_CONNECTION_STRING"); ok {
		return v, nil
	}

	log.Info().Msg("Creating cockroach container...")
	sqlConnectionString, _, err := testcontainers.CreatePostgresServerContainer(t)
	if err != nil {
		return "", errors.Wrap(err, "failed to create cockroach container")
	}

	return sqlConnectionString, nil
}

func getDBPool(sqlConnectionString string) (*pgxpool.Pool, error) {
	pgxPoolConfig, err := pgxpool.ParseConfig(sqlConnectionString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse pgxpool config")
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pgxpool")
	}

	err = dbPool.Ping(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "failed to ping pgxpool")
	}

	return dbPool, nil
}

func createDB(srcDBPool *pgxpool.Pool, migrationProps *MigrationProperties) (*pgxpool.Pool, error) {
	dbName := strings.ReplaceAll("test"+uuid.NewString(), "-", "")

	_, err := srcDBPool.Exec(context.Background(), "CREATE DATABASE "+dbName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create database")
	}

	u, err := url.Parse(srcDBPool.Config().ConnString())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse connection string")
	}

	u.Path = dbName
	sqlConnectionString := u.String()

	dbPool, err := getDBPool(sqlConnectionString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create db pool")
	}

	err = migrationProps.ApplyMigrations(sqlConnectionString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to apply migrations")
	}

	return dbPool, nil
}
