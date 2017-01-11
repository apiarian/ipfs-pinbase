package handlers

import (
	"fmt"
	"net/http"
)

type Login struct {
}

func NewLogin() *Login {
	return &Login{}
}

func (l *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "not implemented yet...")
}
