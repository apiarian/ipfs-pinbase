//go:generate goagen bootstrap -d github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase/design

package main

import (
	"log"
	"time"

	"github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase/app"
	"github.com/apiarian/ipfs-pinbase/pinbase"
	"github.com/apiarian/ipfs-pinbase/pinbase/bolt"
	"github.com/apiarian/ipfs-pinbase/pinbase/ipfs"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
)

func main() {
	P := bolt.NewClient("pinbase.db")
	if err := P.Open(); err != nil {
		log.Fatal("failed to open database connection:", err)
	}
	defer P.Close()

	I, err := ipfs.NewIPFSClient("127.0.0.1:5001")
	if err != nil {
		log.Fatalf("failed to create IPFS client")
	}
	go func(i *ipfs.IPFSClient) {
		for {
			if !i.Ping() {
				log.Print("failed to ping ipfs node")
			}
			time.Sleep(3 * time.Second)
		}
	}(I)

	done := make(chan struct{})

	go pinbase.ManagePins(done, P.PinBackend(), I, 5*time.Second)

	// Create service
	service := goa.New("pinbase")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "party" controller
	c := NewPartyController(service, P)
	app.MountPartyController(service, c)
	// Mount "pin" controller
	c2 := NewPinController(service, P)
	app.MountPinController(service, c2)

	// Start service
	if err := service.ListenAndServe(":3000"); err != nil {
		service.LogError("startup", "err", err)
	}

	close(done)
	log.Print("done")
}
