package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	conn, err := sql.Open("pgx", "host=localhost user=postgres dbname=bluehorse password= postgres sslmode=disable port=5432")
	if err != nil {
		return nil, fmt.Errorf("error connecting to postgres db")
	}

	return conn, nil
}

func MigrateFs(DB *sql.DB, MigrateFS fs.FS, Dir string) error {
	goose.SetBaseFS(MigrateFS)

	defer func() {
		goose.SetBaseFS(nil)
	}()

	return Migrate(DB, Dir)

}

func Migrate(DB *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		log.Fatal(err)
	}

	err = goose.Up(DB, dir)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
