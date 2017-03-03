package main

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/ipfs/go-ipfs-api"
)

func TestPinningMultipleLevels(t *testing.T) {
	s, err := newShellForNode(0)
	if err != nil {
		t.Fatalf("failed to create shell: %+v", err)
	}

	pins, err := s.Pins()
	if err != nil {
		t.Errorf("failed to get initial pins: %+v", err)
	}

	for h, i := range pins {
		if i.Type == shell.DirectPin || i.Type == shell.RecursivePin {
			err = s.Unpin(h)
			if err != nil {
				t.Errorf("failed to unpin %s: %+v", h, err)
			}
		}
	}

	hObjRoot, err := s.AddNoPin(bytes.NewBufferString("root object"))
	if err != nil {
		t.Errorf("failed to create root object: %+v", err)
	}

	hObjLevel1, err := s.AddNoPin(bytes.NewBufferString("level 1"))
	if err != nil {
		t.Errorf("failed to create first level object: %+v", err)
	}

	hObjLevel2, err := s.AddNoPin(bytes.NewBufferString("level 2"))
	if err != nil {
		t.Errorf("failed to create second level object: %+v", err)
	}

	hObjLevel1, err = s.PatchLink(hObjLevel1, "level2", hObjLevel2, false)
	if err != nil {
		t.Errorf("failed to link level 2 to level 1: %+v", err)
	}

	hObjRoot, err = s.PatchLink(hObjRoot, "level1", hObjLevel1, false)
	if err != nil {
		t.Errorf("failed to link level 1 to root: %+v", err)
	}

	objRoot, err := s.ObjectGet(hObjRoot)
	if err != nil {
		t.Errorf("failed to get root object: %+v", err)
	}

	objLevel1, err := s.ObjectGet(hObjLevel1)
	if err != nil {
		t.Errorf("failed to get level 1 object: %+v", err)
	}

	objLevel2, err := s.ObjectGet(hObjLevel2)
	if err != nil {
		t.Errorf("failed to get level 2 object: %+v", err)
	}

	if len(objRoot.Links) != 1 {
		t.Errorf("root should have one link: %+v", objRoot.Links)
	}

	if len(objLevel1.Links) != 1 {
		t.Errorf("first-level object should have one link: %+v", objLevel1.Links)
	}

	if len(objLevel2.Links) != 0 {
		t.Errorf("second-level object should have no links: %+v", objLevel2.Links)
	}

	t.Logf("root (%s): %+v", hObjRoot, objRoot)
	t.Logf("lvl1 (%s): %+v", hObjLevel1, objLevel1)
	t.Logf("lvl2 (%s): %+v", hObjLevel2, objLevel2)

	checkPins := func(tag string, want map[string]shell.PinInfo) {
		got, err := s.Pins()
		if err != nil {
			t.Errorf("%s: failed to get pins: %+v", tag, err)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("%s:\nwant: %+v,\n got: %+v", tag, want, got)
		}
	}

	pinAndCheck := func(tag, h string, want map[string]shell.PinInfo) {
		err := s.Pin(h)
		if err != nil {
			t.Errorf("%s: failed to pin %s: %+v", tag, h, err)
			return
		}

		time.Sleep(100 * time.Millisecond)

		checkPins(tag, want)
	}

	unpinAndCheck := func(tag, h string, want map[string]shell.PinInfo) {
		err := s.Unpin(h)
		if err != nil {
			t.Errorf("%s: failed to unpin %s: %+v", tag, h, err)
			return
		}

		time.Sleep(100 * time.Millisecond)

		checkPins(tag, want)
	}

	checkPins("initial setup", map[string]shell.PinInfo{})

	// Pin an intermediate thing, then the root

	pinAndCheck(
		"add lvl1",
		hObjLevel1,
		map[string]shell.PinInfo{
			hObjLevel1: shell.PinInfo{Type: shell.RecursivePin},
			hObjLevel2: shell.PinInfo{Type: shell.IndirectPin},
		},
	)

	pinAndCheck(
		"add root",
		hObjRoot,
		map[string]shell.PinInfo{
			hObjRoot:   shell.PinInfo{Type: shell.RecursivePin},
			hObjLevel1: shell.PinInfo{Type: shell.RecursivePin},
			hObjLevel2: shell.PinInfo{Type: shell.IndirectPin},
		},
	)

	// Unpin the intermeidate thing

	unpinAndCheck(
		"unpin lvl1",
		hObjLevel1,
		map[string]shell.PinInfo{
			hObjRoot:   shell.PinInfo{Type: shell.RecursivePin},
			hObjLevel1: shell.PinInfo{Type: shell.IndirectPin},
			hObjLevel2: shell.PinInfo{Type: shell.IndirectPin},
		},
	)

	// Pin the intermediate thing again

	pinAndCheck(
		"add lvl1",
		hObjLevel1,
		map[string]shell.PinInfo{
			hObjRoot:   shell.PinInfo{Type: shell.RecursivePin},
			hObjLevel1: shell.PinInfo{Type: shell.RecursivePin},
			hObjLevel2: shell.PinInfo{Type: shell.IndirectPin},
		},
	)

	// Unpin the root thing

	unpinAndCheck(
		"unpin root",
		hObjRoot,
		map[string]shell.PinInfo{
			hObjLevel1: shell.PinInfo{Type: shell.RecursivePin},
			hObjLevel2: shell.PinInfo{Type: shell.IndirectPin},
		},
	)

	// Unpin the intermediate thing

	unpinAndCheck(
		"unpin lvl1",
		hObjLevel1,
		map[string]shell.PinInfo{},
	)
}

func TestPinningRemoteThings(t *testing.T) {
	s0, err := newShellForNode(0)
	if err != nil {
		t.Fatalf("failed to create shell 0: %+v", err)
	}

	s1, err := newShellForNode(1)
	if err != nil {
		t.Fatalf("failed to create shell 1: %+v", err)
	}

	hObj, err := s0.Add(bytes.NewBufferString("this is a test of remote things"))
	if err != nil {
		t.Errorf("failed to create a thing on node 0: %+v", err)
	}

	pins, err := s1.Pins()
	if err != nil {
		t.Errorf("did not get pins: %+v", err)
	}

	_, pinned := pins[hObj]
	if pinned {
		t.Errorf("apparently the remote thing is already pinned")
	}

	err = s1.Pin(hObj)
	if err != nil {
		t.Errorf("failed to pin the remote thing: %+v", err)
	}

	time.Sleep(100 * time.Millisecond)

	pins, err = s1.Pins()
	if err != nil {
		t.Errorf("failed to get pins again: %+v", err)
	}

	_, pinned = pins[hObj]
	if !pinned {
		t.Errorf("apparently failed to actually pin the thing: %+v", err)
	}
}
