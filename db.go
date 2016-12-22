package main

import (
	"database/sql"
	"io"
	"log"
	"os"

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

func migrate(db *sql.DB, M []migrations, start int) error {
	if start == 0 {
		_, err := db.Exec(`
			create table versions (
				number integer primary key not null,
				applied_on text not null default (strftime('%Y-%m-%d %H:%M:%f', 'now'))
			);
		`)
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

		_, err := db.Exec("insert into versions (number) values (?)", i)
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

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.Wrap(err, "open db")
	}

	if !dbExists {
		err = migrate(db, migrations, 0)
		if err != nil {
			return nil, errors.Wrap(err, "migrate initial db")
		}
	} else {
		r, err := db.QueryRow(`select max(number) from versions`)
		if err != nil {
			return nil, errors.Wrap(err, "find latest version")
		}

		var n int
		err = r.Scan(&n)
		if err != nil {
			return nil, errors.Wrap(err, "extract latest version number")
		}

		if n < len(migrations)-1 {
			db.Close()

			err := copyFile(dbPath+".bkp", dbPath)
			if err != nil {
				return nil, errors.Wrap(err, "backup db")
			}

			db, err = sql.Open("sqlite3", dbPath)
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

	_, err := io.Copy(out, in)
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
