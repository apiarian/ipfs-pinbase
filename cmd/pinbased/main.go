package main

import (
	"net/http"

	"goji.io"
)

func main() {
	root := goji.NewMux()

	http.ListenAndServe("localhost:3000", root)
}
