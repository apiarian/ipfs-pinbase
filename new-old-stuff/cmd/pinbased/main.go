package main

import (
	"log"
	"net/http"

	"goji.io"
	"goji.io/pat"

	"github.com/apiarian/ipfs-pinbase/handlers"
)

func main() {
	root := goji.NewMux()

	login, err := handlers.NewLogin([]byte("foo"), nil)
	if err != nil {
		log.Fatalf("failed to create login handler: %+v", err)
	}
	root.Handle(pat.Post("/login"), login)

	apiAddress := "localhost:3000"
	log.Print("listening on ", apiAddress)
	http.ListenAndServe(apiAddress, root)
}
