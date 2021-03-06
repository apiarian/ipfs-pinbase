// Code generated by goagen v1.1.0-dirty, command line:
// $ goagen
// --design=github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase/design
// --out=$(GOPATH)/src/github.com/apiarian/ipfs-pinbase/cmd/ipfs-pinbase
// --version=v1.1.0-dirty
//
// API "pinbase": Application Resource Href Factories
//
// The content of this file is auto-generated, DO NOT MODIFY

package app

import (
	"fmt"
	"strings"
)

// PartyHref returns the resource href.
func PartyHref(partyHash interface{}) string {
	parampartyHash := strings.TrimLeftFunc(fmt.Sprintf("%v", partyHash), func(r rune) bool { return r == '/' })
	return fmt.Sprintf("/api/parties/%v", parampartyHash)
}

// PinHref returns the resource href.
func PinHref(partyHash, pinHash interface{}) string {
	parampartyHash := strings.TrimLeftFunc(fmt.Sprintf("%v", partyHash), func(r rune) bool { return r == '/' })
	parampinHash := strings.TrimLeftFunc(fmt.Sprintf("%v", pinHash), func(r rune) bool { return r == '/' })
	return fmt.Sprintf("/api/parties/%v/pins/%v", parampartyHash, parampinHash)
}
