package test

type MigratorFunc func(string) error

type MigrationProperties struct {
	ApplyMigrations MigratorFunc
}

type MigrationProperty func(*MigrationProperties)

func WithMigrator(f MigratorFunc) MigrationProperty {
	return func(props *MigrationProperties) {
		props.ApplyMigrations = f
	}
}
