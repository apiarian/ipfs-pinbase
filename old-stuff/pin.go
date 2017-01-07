package main

import (
	"log"

	"github.com/apiarian/go-ipfs-api"
	"github.com/pkg/errors"
)

type PinIntention struct {
	PartyHash  string
	ObjectHash string
	Interested bool
}

func ManagePins(
	s *shell.Shell,
	iv []PinIntention,
	c <-chan PinIntention,
	done <-chan struct{},
) {
	pinTrack := make(map[string]map[string]bool)

	for _, v := range iv {
		if _, ok := pinTrack[v.ObjectHash]; !ok {
			pinTrack[v.ObjectHash] = make(map[string]bool)
		}

		pinTrack[v.ObjectHash][v.PartyHash] = v.Interested
	}

	err := reconcilePinNeeds(s, overallPinInterests(pinTrack))
	if err != nil {
		log.Fatalf("failed to initialize pin state: %+v")
	}

	for {
		select {
		case v := <-c:
			if _, ok := pinTrack[v.ObjectHash]; !ok {
				pinTrack[v.ObjectHash] = make(map[string]bool)
			}

			original, existed := pinTrack[v.ObjectHash][v.PartyHash]

			pinTrack[v.ObjectHash][v.PartyHash] = v.Interested

			if !existed || (existed && original != v.Interested) {
				err := reconcilePinNeeds(s, overallPinInterests(pinTrack))
				if err != nil {
					log.Printf("failed to update pin state: %+v", err)
				}
			}
		case <-done:
			return
		}
	}
}

func overallPinInterests(track map[string]map[string]bool) map[string]bool {
	interests := make(map[string]bool)

	for object, parties := range track {
		var interest bool

	ObjectLoop:
		for _, interested := range parties {
			if interested {
				interest = true
				break ObjectLoop
			}
		}

		interests[object] = interest
	}

	return interests
}

func reconcilePinNeeds(s *shell.Shell, pins map[string]bool) error {
	cur, err := s.Pins()
	if err != nil {
		return errors.Wrap(err, "failed to get current pins state")
	}

	errs := make([]error, 0)

	for hash, interested := range pins {
		_, pinned := cur[hash]

		if interested && !pinned {
			err = s.Pin(hash)
			if err != nil {
				errs = append(
					errs,
					errors.Wrapf(err, "failed to pin hash %s", hash),
				)
			}
		}

		if !interested && pinned {
			err = s.Unpin(hash)
			if err != nil {
				errs = append(
					errs,
					errors.Wrapf(err, "failed to unpin hash %s", hash),
				)
			}
		}
	}

	if len(errs) > 0 {
		return errors.Errorf("failed to update some pin states: %+v", errs)
	}

	return nil
}
