package ipfs

import (
	"strings"

	"github.com/apiarian/ipfs-pinbase/pinbase"
	"github.com/ipfs/go-ipfs-api"
	"github.com/pkg/errors"
)

type IPFSClient struct {
	s *shell.Shell
}

func NewIPFSClient(apiAddr string) (*IPFSClient, error) {
	s := shell.NewShell(apiAddr)

	return &IPFSClient{
		s: s,
	}, nil
}

func (ic *IPFSClient) Ping() bool {
	return ic.s.IsUp()
}

func (ic *IPFSClient) Pin(h pinbase.Hash) error {
	return errors.Wrap(ic.s.Pin("/ipfs/"+string(h)), "pin hash")
}

func (ic *IPFSClient) Unpin(h pinbase.Hash) error {
	err := ic.s.Unpin("/ipfs/" + string(h))

	if err != nil && strings.HasSuffix(err.Error(), "not pinned") {
		return nil
	}

	return errors.Wrap(err, "unpin hash")
}

func (ic *IPFSClient) Pins() (map[pinbase.Hash]struct{}, error) {
	pins, err := ic.s.Pins()
	if err != nil {
		return nil, errors.Wrap(err, "get pins")
	}

	r := make(map[pinbase.Hash]struct{})
	for h, t := range pins {
		if t.Type == shell.RecursivePin || t.Type == shell.DirectPin {
			r[pinbase.Hash(h)] = struct{}{}
		}
	}

	return r, nil
}

var _ pinbase.PinJuggler = &IPFSClient{}
