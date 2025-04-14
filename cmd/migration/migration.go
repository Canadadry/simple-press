package migration

import (
	"app/config"
	"app/pkg/dbconn"
	"app/pkg/migration"
	"embed"
	"fmt"
	"io/fs"
)

//go:embed migrations
var migrations embed.FS

const Action = "migration"

func Run(c config.Parameters) error {
	db, err := dbconn.Open(c.DatabaseUrl)
	if err != nil {
		return fmt.Errorf("cannot open database %w", err)
	}
	defer db.Close()
	sub, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("cannot read migration files %w", err)
	}
	readdirSub, ok := sub.(fs.ReadDirFS)
	if !ok {
		return fmt.Errorf("cannot walk sub dir 'migration'")
	}
	err = migration.Migrate(db, "mysql", readdirSub)
	if err != nil {
		return fmt.Errorf("cannot run migration %w", err)
	}
	return nil
}
