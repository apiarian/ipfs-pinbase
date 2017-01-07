package main

import (
	"flag"
	"log"
	"os/user"
	"path/filepath"

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
	_ = pins
}
