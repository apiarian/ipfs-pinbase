package main

import (
	"log"

	"github.com/apiarian/go-ipfs-api"
)

func main() {
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
