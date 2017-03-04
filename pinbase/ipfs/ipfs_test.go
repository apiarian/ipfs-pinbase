package ipfs

import (
	"bytes"
	"testing"
	"time"

	"github.com/apiarian/ipfs-pinbase/pinbase"
)

func TestPinning(t *testing.T) {
	s0, err := newShellForNode(0)
	if err != nil {
		t.Fatalf("failed to get shell: %+v", err)
	}

	apiAddr, err := addressForNode(1)
	if err != nil {
		t.Fatalf("failed to get node address: %+v", err)
	}

	c, err := NewIPFSClient(apiAddr)
	if err != nil {
		t.Fatalf("failed to get client: %+v", err)
	}

	if !c.Ping() {
		t.Fatalf("failed to ping client")
	}

	h1, err := s0.Add(bytes.NewBufferString("a thing"))
	if err != nil {
		t.Errorf("failed to create object 1: %+v", err)
	}

	h2, err := s0.Add(bytes.NewBufferString("another thing"))
	if err != nil {
		t.Errorf("failed to create object 2: %+v", err)
	}

	pins, err := c.Pins()
	if err != nil {
		t.Errorf("failed to get pins: %+v", err)
	}

	if _, pinned := pins[pinbase.Hash(h1)]; pinned {
		t.Errorf("object 1 (%s) somehow pinned already: %+v", h1, pins)
	}

	if _, pinned := pins[pinbase.Hash(h2)]; pinned {
		t.Errorf("object 1 (%s) somehow pinned already: %+v", h2, pins)
	}

	err = c.Pin(pinbase.Hash(h1))
	if err != nil {
		t.Errorf("failed to pin object 1: %+v", err)
	}

	err = c.Pin(pinbase.Hash(h2))
	if err != nil {
		t.Errorf("failed to pin object 2: %+v", err)
	}

	time.Sleep(10 * time.Millisecond)

	pins, err = c.Pins()
	if err != nil {
		t.Errorf("failed to get pins")
	}

	if _, pinned := pins[pinbase.Hash(h1)]; !pinned {
		t.Errorf("object 1 (%s) not pinned: %+v", h1, pins)
	}

	if _, pinned := pins[pinbase.Hash(h2)]; !pinned {
		t.Errorf("object 1 (%s) not pinned: %+v", h2, pins)
	}

	err = c.Unpin(pinbase.Hash(h1))
	if err != nil {
		t.Errorf("failed to unpin object 1: %+v", err)
	}

	time.Sleep(10 * time.Millisecond)

	pins, err = c.Pins()
	if err != nil {
		t.Errorf("failed to get pins")
	}

	if _, pinned := pins[pinbase.Hash(h1)]; pinned {
		t.Errorf("object 1 (%s) still pinned: %+v", h1, pins)
	}

	if _, pinned := pins[pinbase.Hash(h2)]; !pinned {
		t.Errorf("object 1 (%s) not pinned: %+v", h2, pins)
	}

	err = c.Unpin(pinbase.Hash(h2 + "foobar"))
	if err == nil {
		t.Errorf("did not get an error unpinning junk")
	}
}
