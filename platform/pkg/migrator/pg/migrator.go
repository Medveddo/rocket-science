package pg

import (
	"database/sql"

	"github.com/pressly/goose/v3"

	"github.com/Medveddo/rocket-science/platform/pkg/migrator"
)

type Migrator struct {
	db            *sql.DB
	migrationsDir string
}

func NewMigrator(db *sql.DB, migrationsDir string) migrator.Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

func (m *Migrator) Up() error {
	return goose.Up(m.db, m.migrationsDir)
}

func (m *Migrator) Down() error {
	return goose.Down(m.db, m.migrationsDir)
}

func (m *Migrator) Status() error {
	return goose.Status(m.db, m.migrationsDir)
}

func (m *Migrator) Version() error {
	return goose.Version(m.db, m.migrationsDir)
}
