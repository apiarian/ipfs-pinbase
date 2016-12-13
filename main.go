package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/apiarian/go-ipfs-api"
)

func main() {
	usr, _ := user.Current()

	var dbPath, ipfsAPI string

	flag.StringVar(
		&dbPath,
		"db-path",
		filepath.Join(usr.HomeDir, ".pinbase.db"),
		"path to the sqlite3 pin database",
	)
	flag.StringVar(
		&ipfsAPI,
		"ipfs-api-address",
		"localhost:5001",
		"address of the IPFS API",
	)
	flag.Parse()

	s := shell.NewShell(ipfsAPI)
	if !s.IsUp() {
		log.Fatal("IPFS API is not reachable")
	}

	pins, err := s.Pins()
	if err != nil {
		log.Fatalf("failed to get pins: %+v", err)
	}

	dbExists, err := exists(dbPath)
	if err != nil {
		log.Fatalf("failed to check the existence of the database: %+v", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("failed to open DB: %+v", err)
	}
	defer db.Close()

	if !dbExists {
		log.Printf("creating a new database")

		schema := `
		create table pins (
			id integer primary key not null,
			hash text not null,
			created_on text not null default (strftime('%Y-%m-%d %H:%M:%f', 'now')),
			description text not null
		);
		create table parties (
			id integer primary key not null,
			hash text not null,
			kind text not null,
			description text not null
		);
		create table pin_parties (
			id integer primary key not null,
			pin_id integer not null,
			party_id integer not null,
			foreign key(pin_id) references pins(id),
			foreign key(party_id) references parties(id)
		);
		`
		_, err = db.Exec(schema)
		if err != nil {
			log.Printf("error: %s with %s\n", err, schema)
			return
		}
	}

	insertPartySmt, err := db.Prepare("insert into parties (hash, kind, description) values (?, ?, ?)")
	if err != nil {
		log.Fatalf("failed to prepare insert party statement: %+v", err)
	}

	insertPinSmt, err := db.Prepare("insert into pins (hash, created_on, description) values (?, ?, ?)")
	if err != nil {
		log.Fatalf("failed to prepare insert pin statement: %+v", err)
	}

	insertPinPartySmt, err := db.Prepare("insert into pin_parties (pin_id, party_id) values (?, ?)")
	if err != nil {
		log.Fatalf("failed to prepare insert pin-party statement: %+v", err)
	}

	if !dbExists {
		log.Printf("initializing the database")

		fullID, err := s.ID()
		if err != nil {
			log.Fatalf("failed to get shell ID: %+v", err)
		}

		r, err := insertPartySmt.Exec(fullID.ID, "node", "this node")
		if err != nil {
			log.Fatalf("failed to insert node party into database: %+v", err)
		}
		party_id, err := r.LastInsertId()
		if err != nil {
			log.Fatalf("failed to get id of inserted node party: %+v", err)
		}

		for h, i := range pins {
			if i.Type == shell.IndirectPin {
				continue
			}

			r, err = insertPinSmt.Exec(h, time.Now().UTC(), "unknown")
			if err != nil {
				log.Fatalf("failed to insert pin into database: %+v", err)
			}
			pin_id, err := r.LastInsertId()
			if err != nil {
				log.Fatalf("failed to get id of inserted pin: %+v", err)
			}

			_, err = insertPinPartySmt.Exec(pin_id, party_id)
			if err != nil {
				log.Fatalf("failed to insert pin-party into database: %+v", err)
			}
		}
	}
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
