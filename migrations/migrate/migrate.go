package migrate

import (
	"database/sql"
	"embed"
	"github.com/pressly/goose/v3"
	"io/fs"
)

var Migrations embed.FS

func Migrate(dsn string, path fs.FS) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	goose.SetBaseFS(path)
	return goose.Up(db, "migrations/20240121083325_person_table.sql")
}
