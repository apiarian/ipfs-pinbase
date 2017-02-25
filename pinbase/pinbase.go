package pinbase

import (
	"time"

	"github.com/pkg/errors"
)

type Hash string

type PartyCreate struct {
	ID          Hash
	Description string
}

type PartyEdit struct {
	Description string
}

type PartyView struct {
	ID          Hash
	Description string
}

func (pv *PartyView) String() string {
	return string(pv.ID) + ": " + pv.Description
}

type PinCreate struct {
	ID         Hash
	Aliases    []string
	WantPinned bool
}

type PinEdit struct {
	Aliases    []string
	WantPinned bool
}

type PinView struct {
	ID         Hash
	Aliases    []string
	WantPinned bool
	Status     PinStatus
	LastError  error
}

type PinStatus int

const (
	PinPending PinStatus = iota
	PinPinned
	PinUnpinned
	PinError
	PinFatal
	numPinStatuses
)

type PinService interface {
	Parties() ([]*PartyView, error)
	Party(Hash) (*PartyView, error)

	CreateParty(*PartyCreate) error
	DeleteParty(Hash) error
	UpdateParty(Hash, *PartyEdit) error

	Pins(partyID Hash) ([]*PinView, error)
	Pin(partyID, pinID Hash) (*PinView, error)

	CreatePin(partyID Hash, pc *PinCreate) error
	DeletePin(partyID, pinID Hash) error
	UpdatePin(partyID Hash, pe *PinEdit) error
}

type PinBackend interface {
	PinProcessorBump() <-chan struct{}
	PinRequirements() map[Hash]bool
	NotifyPin(pinID Hash, s *PinBackendState)
}

type PinBackendState struct {
	Status    PinStatus
	LastError error
}

type PinJuggler interface {
	Pin(Hash) error
	Unpin(Hash) error
	Pins() (map[Hash]struct{}, error)
}

func ManagePins(
	done <-chan struct{},
	pb PinBackend,
	pj PinJuggler,
	maxInterval time.Duration,
) {
	processPins(pb, pj)

	t := time.NewTimer(maxInterval)

	for {
		select {
		case <-pb.PinProcessorBump():
			processPins(pb, pj)

		case <-t.C:
			processPins(pb, pj)

		case <-done:
			return
		}

		if !t.Stop() {
			<-t.C
		}
		t.Reset(maxInterval)
	}
}

func processPins(pb PinBackend, pj PinJuggler) {
	pr := pb.PinRequirements()

	ps, err := pj.Pins()

	if err != nil {
		for h, _ := range pr {
			pb.NotifyPin(
				h,
				&PinBackendState{
					Status:    PinError,
					LastError: errors.Wrap(err, "get initial pins"),
				},
			)
		}

		return
	}

	for h, want := range pr {
		_, pinned := ps[h]

		var pbs PinBackendState

		switch {
		case want && pinned:
			pbs = PinBackendState{PinPinned, nil}

		case want && !pinned:
			err = pj.Pin(h)
			if err != nil {
				pbs = PinBackendState{PinError, errors.Wrap(err, "pinning unpinned pin")}
			} else {
				pbs = PinBackendState{PinPinned, nil}
			}

		case !want && pinned:
			err = pj.Unpin(h)
			if err != nil {
				pbs = PinBackendState{PinError, errors.Wrap(err, "unpinning pinned pin")}
			} else {
				pbs = PinBackendState{PinUnpinned, nil}
			}

		case !want && !pinned:
			pbs = PinBackendState{PinUnpinned, nil}

		default:
			panic("somehow failed to account for the combinations of 2 booleans")
		}

		pb.NotifyPin(h, &pbs)
	}
}
