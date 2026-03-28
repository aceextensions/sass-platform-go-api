package migrations

import (
	"embed"
)

//go:embed *.sql
var FS embed.FS

func GetMigrationFiles() (embed.FS, error) {
	return FS, nil
}
