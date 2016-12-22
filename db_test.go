package main

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigrate(t *testing.T) {
	db, err := sql.Open("sqlite3", "file::memory:?_loc=UTC")
	fatalIfErr(t, "failed to create db", err)
	defer db.Close()

	err = migrate(db, nil, initialMigrationNumber)
	fatalIfErr(t, "failed to apply a bare migration", err)

	rows, err := db.Query("select * from versions order by applied_on;")
	fatalIfErr(t, "failed to get all versions in the db", err)
	defer rows.Close()

	for rows.Next() {
		var appliedOn time.Time
		var number int
		err := rows.Scan(&appliedOn, &number)
		fatalIfErr(t, "failed to scan the row", err)
	}
	t.Error("bump")
}
