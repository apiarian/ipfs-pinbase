package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"

	"github.com/pkg/errors"
)

type migration struct {
	description string
	statement   string
}

var migrations []migration = []migration{
	migration{
		description: "create parties table",
		statement: `
			create table parties (
				hash text primary key not null,
				kind text not null,
				description text not null
			);`,
	},
}

const initialMigrationNumber = -1

func migrate(db *sqlx.DB, M []migration, start int) error {
	if start == initialMigrationNumber {
		_, err := db.Exec(`
			create table schema_versions (
				number integer primary key not null,
				applied_on timestamp not null default (strftime('%Y-%m-%d %H:%M:%f', 'now'))
			);
			insert into schema_versions (number) values (?);
		`, initialMigrationNumber)
		if err != nil {
			return errors.Wrapf(err, "creating schema_versions table")
		}
	}

	log.Printf("migrating database")

	for i, m := range M {
		if i < start {
			continue
		}

		log.Printf("  applying %d: %s", i, m.description)

		_, err := db.Exec(m.statement)
		if err != nil {
			return errors.Wrapf(err, "applying migration %d (%s)", i, m.description)
		}

		_, err = db.Exec("insert into schema_versions (number) values (?)", i)
		if err != nil {
			return errors.Wrapf(err, "inserting version number %d", i)
		}
	}

	log.Printf("migration complete")

	return nil
}

func pathToDSN(path string) string {
	return fmt.Sprintf("file:%s?_loc=UTC", path)
}

func Open(path string) (*sqlx.DB, error) {
	return open(path, migrations)
}

func open(path string, M []migration) (*sqlx.DB, error) {
	var dbExists bool
	var err error

	if path == ":memory:" {
		log.Print("warning: :memory: db should only be used in testing")
	} else {
		dbExists, err = exists(path)
		if err != nil {
			return nil, errors.Wrap(err, "check existence of database")
		}
	}

	db, err := sqlx.Connect("sqlite3", pathToDSN(path))
	if err != nil {
		return nil, errors.Wrap(err, "open db")
	}

	if !dbExists {
		err = migrate(db, M, initialMigrationNumber)
		if err != nil {
			return nil, errors.Wrap(err, "migrate initial db")
		}

		return db, nil
	}

	var n int
	err = db.Get(&n, `select max(number) from schema_versions`)
	if err != nil {
		return nil, errors.Wrap(err, "find latest version")
	}

	if n < len(M)-1 {
		err := db.Close()
		if err != nil {
			return nil, errors.Wrap(err, "close the db")
		}

		err = copyFile(path+".migration", path)
		if err != nil {
			return nil, errors.Wrap(err, "create migration db")
		}

		db, err = sqlx.Connect("sqlite3", pathToDSN(path+".migration"))
		if err != nil {
			return nil, errors.Wrap(err, "open migration db")
		}

		err = migrate(db, M, n+1)
		if err != nil {
			return nil, errors.Wrap(err, "migrate existing db")
		}

		err = db.Close()
		if err != nil {
			return nil, errors.Wrap(err, "close the migration db")
		}

		err = os.Rename(path, path+".bkp")
		if err != nil {
			return nil, errors.Wrap(err, "failed to backup current db")
		}

		err = os.Rename(path+".migration", path)
		if err != nil {
			return nil, errors.Wrap(err, "failed to move migrated db to the main path")
		}

		return Open(path)
	}

	return db, nil
}

func copyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return errors.Wrap(err, "open src")
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return errors.Wrap(err, "create dst")
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return errors.Wrap(err, "copy src to dst")
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return true, err
}
