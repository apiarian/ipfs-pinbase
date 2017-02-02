package pinbase

import (
	"github.com/pkg/errors"
)

type Node struct {
	Hash        string
	Description string
	APIURL      string
}

type NodeStorage interface {
	Nodes() ([]*Node, error)
	NodeForHash(hash string) (*Node, error)
	SaveNode(*Node) error
}

type Party struct {
	Hash        string
	Description string
	Pins        []*Pin
}

type Pin struct {
	Hash    string
	Aliases []string
	Pinned  bool
}

type PartyStorage interface {
	Parties() ([]*Party, error)
	PartyForHash(hash string) (*Party, error)
	SaveParty(*Party) error
}

type Pinner interface {
	Pins() (map[string]struct{}, error)
	Pin(string) error
	Unpin(string) error
}

func StartPinManager(done <-chan struct{}, pnr Pinner) error {
	p, err := pnr.Pins()
	if err != nil {
		return errors.Wrap(err, "get initial pins")
	}

	_ = p
	return nil
}
