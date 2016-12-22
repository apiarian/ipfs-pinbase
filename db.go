package main

import (
	"database/sql"
	"io"
	"log"
	"os"

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

func migrate(db *sql.DB, M []migration, start int) error {
	if start == initialMigrationNumber {
		_, err := db.Exec(`
			create table versions (
				applied_on timestamp primary key not null default (strftime('%Y-%m-%d %H:%M:%f', 'now')),
				number integer not null
			);
			insert into versions (number) values (?);
		`, initialMigrationNumber)
		if err != nil {
			return errors.Wrapf(err, "creating versions table")
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

		_, err = db.Exec("insert into versions (number) values (?)", i)
		if err != nil {
			return errors.Wrapf(err, "inserting version number %d", i)
		}
	}

	log.Printf("migration complete")

	return nil
}

func Open(path string) (*sql.DB, error) {
	dbExists, err := exists(path)
	if err != nil {
		return nil, errors.Wrap(err, "check existence of database")
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, errors.Wrap(err, "open db")
	}

	if !dbExists {
		err = migrate(db, migrations, initialMigrationNumber)
		if err != nil {
			return nil, errors.Wrap(err, "migrate initial db")
		}
	} else {
		var n int
		err := db.QueryRow(`select max(number) from versions`).Scan(&n)
		if err != nil {
			return nil, errors.Wrap(err, "find latest version")
		}

		if n < len(migrations)-1 {
			db.Close()

			err := copyFile(path+".bkp", path)
			if err != nil {
				return nil, errors.Wrap(err, "backup db")
			}

			db, err = sql.Open("sqlite3", path)
			if err != nil {
				return nil, errors.Wrap(err, "reopen db")
			}

			err = migrate(db, migrations, n)
			if err != nil {
				return nil, errors.Wrap(err, "migrate existing db")
			}
		}
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
