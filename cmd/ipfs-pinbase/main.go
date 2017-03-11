//go:generate goagen bootstrap -d github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase/design

package main

import (
	"log"

	"github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase/app"
	"github.com/apiarian/ipfs-pinbase/pinbase/bolt"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
)

func main() {
	P := bolt.NewClient("pinbase.db")
	if err := P.Open(); err != nil {
		log.Fatal("failed to open database connection:", err)
	}
	defer P.Close()

	// Create service
	service := goa.New("pinbase")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, false))
	service.Use(middleware.Recover())

	// Mount "party" controller
	c := NewPartyController(service, P)
	app.MountPartyController(service, c)

	// Start service
	if err := service.ListenAndServe(":3000"); err != nil {
		service.LogError("startup", "err", err)
	}
}
