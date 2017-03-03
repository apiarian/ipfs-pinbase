package ipfs

import (
	"github.com/ipfs/go-ipfs-api"
	"github.com/ipfs/ipfs-pinbase/pinbase"
	"github.com/pkg/errors"
)

type IPFSClient struct {
	s *shell.Shell
}

func NewIPFSClient(apiAddr string) (*IPFSClient, error) {
	s := shell.NewShell(apiAddr)
}

func (ic *IPFSClient) Pin(pinbase.Hash) error {
	return errors.New("not implemented")
}

func (ic *IPFSClient) Unpin(pinbase.Hash) error {
	return errors.New("not implemented")
}

func (ic *IPFSClient) Pins() (map[pinbase.Hash]struct{}, error) {
	return make(map[pinbase.Hash]struct{}), errors.New("not implemented")
}

var _ pinbase.PinJuggler = &IPFSClient{}
