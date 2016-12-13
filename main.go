package main

import (
	"flag"
	"log"
	"os/user"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/apiarian/go-ipfs-api"
)

func main() {
	usr, _ := user.Current()

	var (
		dbPath = flag.String(
			"db-path",
			filepath.Join(usr.HomeDir, "pinbase.db"),
			"path to the sqlite3 pin database",
		)
	)

	s := shell.NewShell("localhost:5001")
	if !s.IsUp() {
		log.Fatal("IPFS API is not reachable")
	}

	pins, err := s.Pins()
	if err != nil {
		log.Fatalf("failed to get pins: %+v", err)
	}

	for h, i := range pins {
		if i.Type == shell.IndirectPin {
			continue
		}

		log.Print(h, i)
	}
}
