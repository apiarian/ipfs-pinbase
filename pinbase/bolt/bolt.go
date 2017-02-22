package bolt

import (
	"github.com/apiarian/ipfs-pinbase/pinbase"
	"github.com/pkg/errors"
)

type PinService struct {
}

//
// pinbase.PinService implementation
//

func (ps *PinService) Parties() ([]*pinbase.PartyView, error) {
	return nil, errors.New("not implemented")
}

func (ps *PinService) Party(h pinbase.Hash) (*pinbase.PartyView, error) {
	return nil, errors.New("not implemented")
}

func (ps *PinService) CreateParty(p *pinbase.PartyCreate) error {
	return errors.New("not implemented")
}

func (ps *PinService) DeleteParty(h pinbase.Hash) error {
	return errors.New("not implemented")
}

func (ps *PinService) UpdateParty(p *pinbase.PartyEdit) error {
	return errors.New("not implemented")
}

func (ps *PinService) Pins(partyID pinbase.Hash) ([]*pinbase.PinView, error) {
	return nil, errors.New("not implemented")
}

func (ps *PinService) Pin(partyID, pinID pinbase.Hash) (*pinbase.PinView, error) {
	return nil, errors.New("not implemented")
}

func (ps *PinService) CreatePin(partyID pinbase.Hash, pc *pinbase.PinCreate) error {
	return errors.New("not implemented")
}

func (ps *PinService) DeletePin(partyID, pinID pinbase.Hash) error {
	return errors.New("not implemented")
}

func (ps *PinService) UpdatePin(partyID pinbase.Hash, pe *pinbase.PinEdit) error {
	return errors.New("not implemented")
}

//
// pinbase.PinBackend implementation
//

func (ps *PinService) PinProcessorBump() <-chan struct{} {
	return make(chan struct{})
}

func (ps *PinService) PinRequirements() map[pinbase.Hash]bool {
	return make(map[pinbase.Hash]bool)
}

func (ps *PinService) NotifyPin(pinID pinbase.Hash, s *pinbase.PinBackendState) {
}

var _ pinbase.PinService = &PinService{}
var _ pinbase.PinBackend = &PinService{}
