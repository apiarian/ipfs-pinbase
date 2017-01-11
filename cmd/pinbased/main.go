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

	root.Handle(
		pat.Post("/login"),
		handlers.NewLogin(),
	)

	apiAddress := "localhost:3000"
	log.Print("listening on ", apiAddress)
	http.ListenAndServe(apiAddress, root)
}
