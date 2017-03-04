package pinbase

import (
	"fmt"
	"strings"
	"testing"

	"reflect"

	"github.com/pkg/errors"
)

type MemoryBackendInfo struct {
	WantPinned       bool
	Status           PinStatus
	LastErrorMessage string
}

func (i *MemoryBackendInfo) String() string {
	return fmt.Sprintf(
		"wanted(%t)-status(%s)-err(%v)",
		i.WantPinned,
		i.Status,
		i.LastErrorMessage,
	)
}

type MemoryBackend struct {
	Pins   map[Hash]*MemoryBackendInfo
	Bumper chan struct{}
}

type MemoryJuggler struct {
	P               map[Hash]struct{}
	PinsShouldError bool
}

func NewMemoryBackend() *MemoryBackend {
	return &MemoryBackend{
		Pins:   make(map[Hash]*MemoryBackendInfo),
		Bumper: make(chan struct{}),
	}
}

func (mb *MemoryBackend) PinProcessorBump() <-chan struct{} {
	return mb.Bumper
}

func (mb *MemoryBackend) PinRequirements() map[Hash]bool {
	r := make(map[Hash]bool)

	for h, i := range mb.Pins {
		r[h] = i.WantPinned
	}

	return r
}

func (mb *MemoryBackend) NotifyPin(p Hash, s *PinBackendState) {
	_, ok := mb.Pins[p]
	if ok {
		mb.Pins[p].Status = s.Status
	} else {
		mb.Pins[p] = &MemoryBackendInfo{
			WantPinned: false,
			Status:     s.Status,
		}
	}

	if s.LastError != nil {
		mb.Pins[p].LastErrorMessage = s.LastError.Error()
	} else {
		mb.Pins[p].LastErrorMessage = ""
	}
}

var _ PinBackend = &MemoryBackend{}

func NewMemoryJuggler() *MemoryJuggler {
	return &MemoryJuggler{
		P:               make(map[Hash]struct{}),
		PinsShouldError: false,
	}
}

func (mj *MemoryJuggler) Pin(h Hash) error {
	if strings.HasPrefix(string(h), "bad") {
		return errors.New("cannot pin bad hash")
	}

	mj.P[h] = struct{}{}
	return nil
}

func (mj *MemoryJuggler) Unpin(h Hash) error {
	if strings.HasPrefix(string(h), "bad") {
		return errors.New("cannot unpin bad hash")
	}

	delete(mj.P, h)
	return nil
}

func (mj *MemoryJuggler) Pins() (map[Hash]struct{}, error) {
	if mj.PinsShouldError {
		return nil, errors.New("can't get no pins")
	}

	return mj.P, nil
}

var _ PinJuggler = &MemoryJuggler{}

func TestProcessPins(t *testing.T) {
	pj := NewMemoryJuggler()
	pj.P[Hash("old")] = struct{}{}

	pb := NewMemoryBackend()
	pb.Pins[Hash("wanted")] = &MemoryBackendInfo{
		WantPinned:       true,
		Status:           PinPending,
		LastErrorMessage: "",
	}

	pb.Pins[Hash("notwanted")] = &MemoryBackendInfo{
		WantPinned:       false,
		Status:           PinPending,
		LastErrorMessage: "",
	}

	pb.Pins[Hash("badjunk")] = &MemoryBackendInfo{
		WantPinned:       true,
		Status:           PinPending,
		LastErrorMessage: "",
	}

	processPins(pb, pj)

	if !reflect.DeepEqual(
		pb.Pins,
		map[Hash]*MemoryBackendInfo{
			Hash("wanted"): &MemoryBackendInfo{
				WantPinned:       true,
				Status:           PinPinned,
				LastErrorMessage: "",
			},
			Hash("notwanted"): &MemoryBackendInfo{
				WantPinned:       false,
				Status:           PinUnpinned,
				LastErrorMessage: "",
			},
			Hash("badjunk"): &MemoryBackendInfo{
				WantPinned:       true,
				Status:           PinError,
				LastErrorMessage: "pinning unpinned pin: cannot pin bad hash",
			},
		},
	) {
		t.Errorf("pin backend state incorrect: %+v", pb.Pins)
	}

	if !reflect.DeepEqual(
		pj.P,
		map[Hash]struct{}{
			Hash("old"):    struct{}{},
			Hash("wanted"): struct{}{},
		},
	) {
		t.Errorf("pin storage state incorrect: %+v", pj.P)
	}

	// just to make sure that we don't do any pinning next
	pb.Pins[Hash("wanted")].WantPinned = false

	pj.PinsShouldError = true

	processPins(pb, pj)

	if !reflect.DeepEqual(
		pb.Pins,
		map[Hash]*MemoryBackendInfo{
			Hash("wanted"): &MemoryBackendInfo{
				WantPinned:       false,
				Status:           PinError,
				LastErrorMessage: "get initial pins: can't get no pins",
			},
			Hash("notwanted"): &MemoryBackendInfo{
				WantPinned:       false,
				Status:           PinError,
				LastErrorMessage: "get initial pins: can't get no pins",
			},
			Hash("badjunk"): &MemoryBackendInfo{
				WantPinned:       true,
				Status:           PinError,
				LastErrorMessage: "get initial pins: can't get no pins",
			},
		},
	) {
		t.Errorf("pin backend state incorrect: %+v", pb.Pins)
	}

	if !reflect.DeepEqual(
		pj.P,
		map[Hash]struct{}{
			Hash("old"):    struct{}{},
			Hash("wanted"): struct{}{},
		},
	) {
		t.Errorf("pin storage state incorrect: %+v", pj.P)
	}

	pj.PinsShouldError = false

	processPins(pb, pj)

	if !reflect.DeepEqual(
		pb.Pins,
		map[Hash]*MemoryBackendInfo{
			Hash("wanted"): &MemoryBackendInfo{
				WantPinned:       false,
				Status:           PinUnpinned,
				LastErrorMessage: "",
			},
			Hash("notwanted"): &MemoryBackendInfo{
				WantPinned:       false,
				Status:           PinUnpinned,
				LastErrorMessage: "",
			},
			Hash("badjunk"): &MemoryBackendInfo{
				WantPinned:       true,
				Status:           PinError,
				LastErrorMessage: "pinning unpinned pin: cannot pin bad hash",
			},
		},
	) {
		t.Errorf("pin backend state incorrect: %+v", pb.Pins)
	}

	if !reflect.DeepEqual(
		pj.P,
		map[Hash]struct{}{
			Hash("old"): struct{}{},
		},
	) {
		t.Errorf("pin storage state incorrect: %+v", pj.P)
	}
}
