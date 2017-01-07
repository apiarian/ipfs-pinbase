package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

func TestMigrateNilMigration(t *testing.T) {
	db, err := sqlx.Connect("sqlite3", pathToDSN(":memory:"))
	fatalIfErr(t, "failed to connect to db", err)
	defer db.Close()

	err = migrate(db, nil, initialMigrationNumber)
	fatalIfErr(t, "failed to apply a bare migration", err)

	type VersionRow struct {
		Number    int
		AppliedOn time.Time `db:"applied_on"`
	}

	versionRows := []VersionRow{}
	err = db.Select(&versionRows, "select * from schema_versions order by number;")
	fatalIfErr(t, "failed to select versions from the database", err)

	if len(versionRows) != 1 {
		t.Errorf("did not find one schema_versions row: %+v", versionRows)
	}

	if versionRows[0].Number != initialMigrationNumber {
		t.Errorf(
			"the initial migration does not have the right version number: %d",
			versionRows[0].Number,
		)
	}
}

func TestMigrateMultipleTimes(t *testing.T) {
	db, err := sqlx.Connect("sqlite3", pathToDSN(":memory:"))
	fatalIfErr(t, "failed to connect to db", err)
	defer db.Close()

	M := []migration{
		migration{
			description: "step 1",
			statement: `
				create table foo (
					id integer primary key not null,
					bar text not null
				);
				insert into foo (bar) values ('hello world');
			`,
		},
	}

	err = migrate(db, M, initialMigrationNumber)
	fatalIfErr(t, "failed to apply 1 step migration", err)

	var maxVersion int
	err = db.Get(&maxVersion, "select max(number) from schema_versions")
	fatalIfErr(t, "failed to get max version", err)

	if maxVersion != 0 {
		t.Errorf("the maximum version number is not 0: %d", maxVersion)
	}

	type Foo struct {
		ID  int
		Bar string
	}

	foos := []Foo{}
	err = db.Select(&foos, "select * from foo order by id;")
	fatalIfErr(t, "failed to get the foos", err)

	if len(foos) != 1 {
		t.Errorf("did not get exactly one foo: %+v", foos)
	}

	if foos[0].ID != 1 {
		t.Errorf("the first foo's ID is not 1: %d", foos[0].ID)
	}
	if foos[0].Bar != "hello world" {
		t.Errorf("the first foo's bar is not right: %s", foos[0].Bar)
	}

	M = append(
		M,
		migration{
			description: "step 2",
			statement: `
				alter table foo add column baz real not null default 0.0;
				insert into foo (bar, baz) values ('oh hi', 3.14);
			`,
		},
	)

	err = migrate(db, M, maxVersion+1)
	fatalIfErr(t, "failed to apply step 2 migration", err)

	err = db.Get(&maxVersion, "select max(number) from schema_versions")
	fatalIfErr(t, "failed to get max version again", err)

	if maxVersion != 1 {
		t.Errorf("the maximum version number is not 1: %d", maxVersion)
	}

	type Foo2 struct {
		ID  int
		Bar string
		Baz float64
	}

	foos2 := []Foo2{}
	err = db.Select(&foos2, "select * from foo order by id;")
	fatalIfErr(t, "failed to get the foos2", err)

	if len(foos2) != 2 {
		t.Errorf("did not get exactly two foos2: %+v", foos2)
	}

	if foos2[0].Baz != 0.0 {
		t.Errorf("the first foo does not have the default 0.0 baz: %f", foos2[0].Baz)
	}
	if foos2[1].Baz != 3.14 {
		t.Errorf("the second foo does not have a 3.14 baz: %f", foos2[0].Baz)
	}
}

func TestMigrateFellSwoop(t *testing.T) {
	db, err := sqlx.Connect("sqlite3", pathToDSN(":memory:"))
	fatalIfErr(t, "failed to connect to db", err)
	defer db.Close()

	M := []migration{
		migration{
			description: "step 1",
			statement: `
				create table foo (
					id integer primary key not null,
					bar text not null
				);
				insert into foo (bar) values ('hello world');
			`,
		},
		migration{
			description: "step 2",
			statement: `
				alter table foo add column baz real not null default 0.0;
				insert into foo (bar, baz) values ('oh hi', 3.14);
			`,
		},
	}

	err = migrate(db, M, initialMigrationNumber)
	fatalIfErr(t, "failed to apply 2 step migration", err)

	var maxVersion int
	err = db.Get(&maxVersion, "select max(number) from schema_versions")
	fatalIfErr(t, "failed to get max version", err)

	if maxVersion != 1 {
		t.Errorf("the maximum version number is not 1: %d", maxVersion)
	}

	type Foo struct {
		ID  int
		Bar string
		Baz float64
	}

	foos := []Foo{}
	err = db.Select(&foos, "select * from foo order by id;")
	fatalIfErr(t, "failed to get the foos", err)

	if len(foos) != 2 {
		t.Errorf("did not get exactly two foos: %+v", foos)
	}

	if foos[0].ID != 1 {
		t.Errorf("the first foo's ID is not 1: %d", foos[0].ID)
	}
	if foos[0].Bar != "hello world" {
		t.Errorf("the first foo's bar is not right: %s", foos[0].Bar)
	}
	if foos[0].Baz != 0.0 {
		t.Errorf("the first foo does not have the default 0.0 baz: %f", foos[0].Baz)
	}

	if foos[1].Baz != 3.14 {
		t.Errorf("the second foo does not have a 3.14 baz: %f", foos[0].Baz)
	}
}

func TestOpen(t *testing.T) {
	dir, err := ioutil.TempDir("", "pinbase-test-open")
	fatalIfErr(t, "failed to create temporary directory", err)
	defer os.RemoveAll(dir)

	M := []migration{
		migration{
			description: "step 1",
			statement: `
						create table foos(
							id integer primary key not null,
							bar text not null
						);
					`,
		},
		migration{
			description: "step 2",
			statement:   "alter table foos add column baz real not null default 0.0",
		},
		migration{
			description: "broken step",
			statement:   "foo bar baz asdf",
		},
	}

	dbPath := filepath.Join(dir, "test.db")

	dbExists, err := exists(dbPath)
	fatalIfErr(t, "failed to check db existence before tests", err)
	if dbExists {
		t.Fatal("database inexplicably exists before the tests")
	}

	const noBackupPresentVersion = -2

	sequence := []struct {
		tag                   string
		M                     []migration
		expectOpenError       bool
		maxVersionExpected    int
		dbShouldExist         bool
		backupShouldExist     bool
		backupVersionExpected int
		migrationShouldExist  bool
	}{
		{
			tag:                   "initial bare database",
			M:                     nil,
			expectOpenError:       false,
			maxVersionExpected:    initialMigrationNumber,
			dbShouldExist:         true,
			backupShouldExist:     false,
			backupVersionExpected: noBackupPresentVersion,
			migrationShouldExist:  false,
		},
		{
			tag:                   "first migration",
			M:                     M[0:1],
			expectOpenError:       false,
			maxVersionExpected:    0,
			dbShouldExist:         true,
			backupShouldExist:     true,
			backupVersionExpected: initialMigrationNumber,
			migrationShouldExist:  false,
		},
		{
			tag:                   "second migration",
			M:                     M[0:2],
			expectOpenError:       false,
			maxVersionExpected:    1,
			dbShouldExist:         true,
			backupShouldExist:     true,
			backupVersionExpected: 0,
			migrationShouldExist:  false,
		},
		{
			tag:                   "broken migration",
			M:                     M[0:3],
			expectOpenError:       true,
			maxVersionExpected:    1,
			dbShouldExist:         true,
			backupShouldExist:     true,
			backupVersionExpected: 0, // means that the backup was not touched either
			migrationShouldExist:  true,
		},
	}

	for _, step := range sequence {
		t.Logf("# %s", step.tag)

		db, err := open(dbPath, step.M)
		if !step.expectOpenError {
			if err != nil {
				t.Fatalf("got an unexpected db open error: %+v", err)
			}
		} else {
			if err == nil {
				t.Fatal("did not get an expected open error")
			}

			db, err = sqlx.Connect("sqlite3", pathToDSN(dbPath))
			fatalIfErr(t, "failed to force db connection", err)
		}

		var maxVersion int
		err = db.Get(&maxVersion, "select max(number) from schema_versions;")
		fatalIfErr(t, "failed to get maxVersion from database", err)

		if maxVersion != step.maxVersionExpected {
			t.Errorf("maxVersion is not %d: %d", step.maxVersionExpected, maxVersion)
		}
		dbExists, err := exists(dbPath)
		fatalIfErr(t, "failed to check if db exists on disk", err)
		if dbExists != step.dbShouldExist {
			t.Error("got %t db exists but expected %t", dbExists, step.dbShouldExist)
		}

		backupExists, err := exists(dbPath + ".bkp")
		fatalIfErr(t, "failed to check if backup db exists on disk", err)
		if backupExists != step.backupShouldExist {
			t.Errorf("got %t backup exists but expected %t", backupExists, step.backupShouldExist)
		}

		if step.backupVersionExpected != noBackupPresentVersion {
			backupDb, err := sqlx.Connect("sqlite3", pathToDSN(dbPath+".bkp"))
			fatalIfErr(t, "failed to connect to backup database", err)
			defer backupDb.Close()

			var backupVersion int
			err = backupDb.Get(&backupVersion, "select max(number) from schema_versions;")
			fatalIfErr(t, "failed to get backupVersion from database", err)

			if backupVersion != step.backupVersionExpected {
				t.Errorf("backupVerion is not %d, %d", step.backupVersionExpected, backupVersion)
			}
		}

		migrationExists, err := exists(dbPath + ".migration")
		fatalIfErr(t, "failed to check if migration db exists on disk", err)
		if migrationExists != step.migrationShouldExist {
			t.Errorf("got %t migration exists but expected %t", migrationExists, step.migrationShouldExist)
		}

		err = db.Close()
		fatalIfErr(t, "failed to close the db", err)
	}
}
