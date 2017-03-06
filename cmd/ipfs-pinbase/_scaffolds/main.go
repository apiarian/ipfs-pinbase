//go:generate goagen bootstrap -d github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase/design

package main

import (
	"github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase/_scaffolds/app"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
)

func main() {
	// Create service
	service := goa.New("pinbase")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "login" controller
	c := NewLoginController(service)
	app.MountLoginController(service, c)
	// Mount "node" controller
	c2 := NewNodeController(service)
	app.MountNodeController(service, c2)

	// Start service
	if err := service.ListenAndServe(":3000"); err != nil {
		service.LogError("startup", "err", err)
	}
}
