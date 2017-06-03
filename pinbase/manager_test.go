package pinbase

import (
	"context"
	"sync"
	"testing"
	"time"
)

type NullBackend struct {
	Bumper chan struct{}
	c      int
	m      *sync.Mutex
}

func NewNullBackend() *NullBackend {
	return &NullBackend{
		Bumper: make(chan struct{}),
		m:      &sync.Mutex{},
	}
}

func (nb *NullBackend) PinProcessorBump() <-chan struct{} {
	return nb.Bumper
}

func (nb *NullBackend) PinsCallCount() int {
	nb.m.Lock()
	defer nb.m.Unlock()
	return nb.c
}

func (nb *NullBackend) PinRequirements() map[Hash]bool {
	nb.m.Lock()
	defer nb.m.Unlock()
	nb.c = nb.c + 1

	return make(map[Hash]bool)
}

func (nb *NullBackend) NotifyPin(_ Hash, _ *PinBackendState) {
	return
}

var _ PinBackend = &NullBackend{}

type NullJuggler struct {
}

func NewNullJuggler() *NullJuggler {
	return &NullJuggler{}
}

func (nj *NullJuggler) Pin(Hash) error {
	return nil
}

func (nj *NullJuggler) Unpin(Hash) error {
	return nil
}

func (nj *NullJuggler) Pins() (map[Hash]struct{}, error) {
	return make(map[Hash]struct{}), nil
}

var _ PinJuggler = &NullJuggler{}

func TestManagePins(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	pb := NewNullBackend()
	pj := NewNullJuggler()

	if c := pb.PinsCallCount(); c != 0 {
		t.Errorf("did not start out with a zero reqs call count: %d", c)
	}

	go ManagePins(ctx, pb, pj, 1*time.Second)

	time.Sleep(10 * time.Millisecond)

	if c := pb.PinsCallCount(); c != 1 {
		t.Errorf("reqs call count should be 1: %d", c)
	}

	go func(c chan<- struct{}) { c <- struct{}{} }(pb.Bumper)

	time.Sleep(10 * time.Millisecond)

	if c := pb.PinsCallCount(); c != 2 {
		t.Errorf("reqs call count should be 2: %d", c)
	}

	// wait for the timeout to trip
	time.Sleep(1 * time.Second)

	if c := pb.PinsCallCount(); c != 3 {
		t.Errorf("reqs call count should be 3: %d", c)
	}

	// wait for another second for the timeout to trip
	time.Sleep(1250 * time.Millisecond)

	if c := pb.PinsCallCount(); c != 4 {
		t.Errorf("reqs call count should be 4: %d", c)
	}

	cancel()

	time.Sleep(10 * time.Millisecond)

	go func(c chan<- struct{}) { c <- struct{}{} }(pb.Bumper)

	time.Sleep(10 * time.Millisecond)

	if c := pb.PinsCallCount(); c != 4 {
		t.Errorf("reqs call count should be 4: %d", c)
	}
}
