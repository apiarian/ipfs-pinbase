package pinbase

import (
	"time"

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

type Intention struct {
	Party   string
	Object  string
	WantPin bool
}

type InterestStatus struct {
	PinnerErr   error
	PinStatuses map[string]PinStatus
}

type PinStatus struct {
	Timestamp   time.Time
	LatestError error
	State       PinState
}

type PinState int

const (
	PinPending PinState = iota
	PinPinned
	PinUnpinned
	PinUnstable
	PinFailed
	numPinStates
)

type InterestTracker interface {
	BootstrapInterest([]Intention)
	UpdateInterest(Intention)
	InterestDigest() map[string]bool
	NotifyState(map[string]struct{}, error, map[string]error)
	Status() *InterestStatus
}

func ManagePins(
	done <-chan struct{},
	pnr Pinner,
	trkr InterestTracker,
	intentions <-chan Intention,
	maxInterval time.Duration,
) {
	tryThePins(pnr, trkr)

	t := time.NewTimer(maxInterval)

	for {
		select {
		case i := <-intentions:
			trkr.UpdateInterest(i)
			tryThePins(pnr, trkr)

		case <-t.C:
			tryThePins(pnr, trkr)

		case <-done:
			return
		}

		if !t.Stop() {
			<-t.C
		}
		t.Reset(maxInterval)
	}
}

func tryThePins(pnr Pinner, trkr InterestTracker) {
	errs := make(map[string]error)

	p, err := pnr.Pins()
	if err != nil {
		trkr.NotifyState(p, errors.Wrap(err, "get initial pins"), errs)

		return
	}

	for hash, want := range trkr.InterestDigest() {
		_, pinned := p[hash]

		if want && !pinned {
			errs[hash] = errors.Wrapf(
				pnr.Pin(hash),
				"pin %s", hash,
			)

		} else if !want && pinned {
			errs[hash] = errors.Wrapf(
				pnr.Unpin(hash),
				"unpin %s", hash,
			)
		}
	}

	p, err = pnr.Pins()
	trkr.NotifyState(p, errors.Wrap(err, "get final state"), errs)
}
