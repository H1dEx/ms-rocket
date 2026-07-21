package migrator

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db *sql.DB
	dir string
}

func NewMigrator(db *sql.DB, dir string) *Migrator {
	return &Migrator{
		db: db,
		dir: dir,
	}
}

func (m *Migrator) Up () error {
	err := goose.Up(m.db, m.dir)
	if err != nil {
		return err
	}
	return nil
}