package test

import (
	"reflect"
	"testing"
	"time"

	"github.com/apiarian/ipfs-pinbase/pinbase"
)

func TestPinServiceHappyPath(t *testing.T, ps pinbase.PinService) {
	// empty service is empty
	parties, err := ps.Parties()
	if err != nil {
		t.Errorf("did not get parties: %+v", err)
	}

	if len(parties) != 0 {
		t.Error("did not get zero parties: %+v", parties)
	}

	// make the service not empty
	err = ps.CreateParty(&pinbase.PartyCreate{
		ID:          pinbase.Hash("foo"),
		Description: "hello world",
	})
	if err != nil {
		t.Errorf("did not create party: %+v", err)
	}

	// list the parties
	parties, err = ps.Parties()
	if err != nil {
		t.Errorf("did not get parties: %+v", err)
	}

	if len(parties) != 1 {
		t.Errorf("did not get one party: %+v", parties)
	}

	if party := parties[0]; !reflect.DeepEqual(
		party,
		&pinbase.PartyView{
			ID:          pinbase.Hash("foo"),
			Description: "hello world",
		},
	) {
		t.Errorf("did not get the party we expected: %+v", party)
	}

	// add another party
	err = ps.CreateParty(&pinbase.PartyCreate{
		ID:          pinbase.Hash("bar"),
		Description: "the party is only just beginning",
	})
	if err != nil {
		t.Errorf("did not create second party: %+v", err)
	}

	parties, err = ps.Parties()
	if err != nil {
		t.Errorf("did not get the parties: %+v", err)
	}

	if len(parties) != 2 {
		t.Errorf("did not get two parties: %+v", parties)
	}

	partiesByKey := make(map[pinbase.Hash]*pinbase.PartyView)
	for _, p := range parties {
		partiesByKey[p.ID] = p
	}

	if !reflect.DeepEqual(
		partiesByKey,
		map[pinbase.Hash]*pinbase.PartyView{
			pinbase.Hash("foo"): &pinbase.PartyView{
				ID:          pinbase.Hash("foo"),
				Description: "hello world",
			},
			pinbase.Hash("bar"): &pinbase.PartyView{
				ID:          pinbase.Hash("bar"),
				Description: "the party is only just beginning",
			},
		},
	) {
		t.Errorf("did not get the parties we wanted: %+v", parties)
	}

	// get an existing party
	party, err := ps.Party(pinbase.Hash("foo"))
	if err != nil {
		t.Errorf("did not get the party: %+v", err)
	}

	if !reflect.DeepEqual(
		party,
		&pinbase.PartyView{
			ID:          pinbase.Hash("foo"),
			Description: "hello world",
		},
	) {
		t.Errorf("did not get the expected party: %+v", party)
	}

	// try to get a nonexistent party
	party, err = ps.Party(pinbase.Hash("baz"))
	if err != nil {
		t.Errorf("failed to try to get the party: %+v", err)
	}

	if party != nil {
		t.Errorf("got a nonexistent party: %+v", party)
	}

	// update a party
	err = ps.UpdateParty(
		pinbase.Hash("bar"),
		&pinbase.PartyEdit{
			Description: "this party will be over soon",
		},
	)
	if err != nil {
		t.Errorf("did not update the party: %+v", err)
	}

	party, err = ps.Party(pinbase.Hash("bar"))
	if err != nil {
		t.Errorf("did not get the party: %+v", err)
	}

	if !reflect.DeepEqual(
		party,
		&pinbase.PartyView{
			ID:          pinbase.Hash("bar"),
			Description: "this party will be over soon",
		},
	) {
		t.Errorf("party was not updated correctly: %+v", party)
	}

	// delete a party
	err = ps.DeleteParty(pinbase.Hash("bar"))
	if err != nil {
		t.Errorf("did not delete the party: %+v", err)
	}

	party, err = ps.Party(pinbase.Hash("bar"))
	if err != nil {
		t.Errorf("failed to try to get the party: %+v", err)
	}

	if party != nil {
		t.Errorf("got a deleted party: %+v", party)
	}

	parties, err = ps.Parties()
	if err != nil {
		t.Errorf("failed to get the parties: %+v", err)
	}

	if len(parties) != 1 {
		t.Errorf("we don't have exaclty one party again: %+v", parties)
	}

	if party := parties[0]; party.ID != pinbase.Hash("foo") {
		t.Errorf("we don't stil have the expected party: %+v", party)
	}

	// delete a nonexistent party
	err = ps.DeleteParty(pinbase.Hash("baz"))
	if err != nil {
		t.Errorf("failed to delete nonexistent party: %+v", err)
	}

	party, err = ps.Party(pinbase.Hash("baz"))
	if err != nil {
		t.Errorf("failed to try to get the party: %+v", err)
	}

	if party != nil {
		t.Errorf("somehow created a party while deleting it: %+v", party)
	}

	//
	// now lets look at the pins
	//

	// nothing to see here yet
	pins, err := ps.Pins(pinbase.Hash("foo"))
	if err != nil {
		t.Errorf("failed to try to get the pins: %+v", err)
	}

	if len(pins) != 0 {
		t.Errorf("did not get zero pins: %+v", pins)
	}

	// create a pin
	err = ps.CreatePin(pinbase.Hash("foo"), &pinbase.PinCreate{
		ID:         pinbase.Hash("bar"),
		Aliases:    []string{"really cool", "super rad"},
		WantPinned: true,
	})
	if err != nil {
		t.Errorf("did not create pin: %+v", err)
	}

	// list the pins
	pins, err = ps.Pins(pinbase.Hash("foo"))
	if err != nil {
		t.Errorf("did not get pins: %+v", err)
	}

	if len(pins) != 1 {
		t.Errorf("did not get one pin: %+v", pins)
	}

	if pin := pins[0]; !reflect.DeepEqual(
		pin,
		&pinbase.PinView{
			ID:         pinbase.Hash("bar"),
			Aliases:    []string{"really cool", "super rad"},
			WantPinned: true,
			Status:     pinbase.PinPending,
			LastError:  nil,
		},
	) {
		t.Errorf("did not get the pin we expected: %+v", pin)
	}

	// add another pin
	err = ps.CreatePin(pinbase.Hash("foo"), &pinbase.PinCreate{
		ID:         pinbase.Hash("abc"),
		Aliases:    []string{"something"},
		WantPinned: false,
	})
	if err != nil {
		t.Errorf("did not create pin: %+v", err)
	}

	pins, err = ps.Pins(pinbase.Hash("foo"))
	if err != nil {
		t.Errorf("did not get the pins: %+v", err)
	}

	if len(pins) != 2 {
		t.Errorf("did not get two pins: %+v", pins)
	}

	pinsByKey := make(map[pinbase.Hash]*pinbase.PinView)
	for _, p := range pins {
		pinsByKey[p.ID] = p
	}

	if !reflect.DeepEqual(
		pinsByKey,
		map[pinbase.Hash]*pinbase.PinView{
			pinbase.Hash("bar"): &pinbase.PinView{
				ID:         pinbase.Hash("bar"),
				Aliases:    []string{"really cool", "super rad"},
				WantPinned: true,
				Status:     pinbase.PinPending,
				LastError:  nil,
			},
			pinbase.Hash("abc"): &pinbase.PinView{
				ID:         pinbase.Hash("abc"),
				Aliases:    []string{"something"},
				WantPinned: false,
				Status:     pinbase.PinPending,
				LastError:  nil,
			},
		},
	) {
		t.Errorf("did not get the pins we wanted: %+v", pins)
	}

	// get an existing pin
	pin, err := ps.Pin(pinbase.Hash("foo"), pinbase.Hash("bar"))
	if err != nil {
		t.Errorf("did not get pin: %+v", err)
	}

	if !reflect.DeepEqual(
		pin,
		&pinbase.PinView{
			ID:         pinbase.Hash("bar"),
			Aliases:    []string{"really cool", "super rad"},
			WantPinned: true,
			Status:     pinbase.PinPending,
			LastError:  nil,
		},
	) {
		t.Errorf("did ont get the expected pin: %+v", err)
	}

	// try to get an nonexistent pin
	pin, err = ps.Pin(pinbase.Hash("foo"), pinbase.Hash("baz"))
	if err != nil {
		t.Errorf("failed to try to get the pin: %+v", err)
	}

	if pin != nil {
		t.Errorf("got a nonexistent pin: %+v", pin)
	}

	// update a pin
	err = ps.UpdatePin(
		pinbase.Hash("foo"),
		pinbase.Hash("bar"),
		&pinbase.PinEdit{
			Aliases:    []string{"really rad"},
			WantPinned: false,
		},
	)
	if err != nil {
		t.Errorf("did not update the pin: %+v", err)
	}

	pin, err = ps.Pin(pinbase.Hash("foo"), pinbase.Hash("bar"))
	if err != nil {
		t.Errorf("did not get the pin: %+v", err)
	}

	if !reflect.DeepEqual(
		pin,
		&pinbase.PinView{
			ID:         pinbase.Hash("bar"),
			Aliases:    []string{"really rad"},
			WantPinned: false,
			Status:     pinbase.PinPending,
			LastError:  nil,
		},
	) {
		t.Errorf("pin was not updated correctly: %+v", pin)
	}

	// delete a pin
	err = ps.DeletePin(pinbase.Hash("foo"), pinbase.Hash("abc"))
	if err != nil {
		t.Errorf("did not delete the pin: %+v", err)

		pin, err = ps.Pin(pinbase.Hash("foo"), pinbase.Hash("abc"))
		if err != nil {
			t.Errorf("failed to try to get the pin: %+v", err)
		}

		if pin != nil {
			t.Errorf("got a deleted pin: %+v", pin)
		}

		pins, err = ps.Pins(pinbase.Hash("foo"))
		if err != nil {
			t.Errorf("failed to get the pins: %+v", err)
		}

		if len(pins) != 1 {
			t.Errorf("we don't have exactly one pin again: %+v", pins)
		}

		if pin := pins[0]; pin.ID != pinbase.Hash("bar") {
			t.Errorf("we don't still have the expected pin: %+v", pin)
		}

		// delte a nonexistent pin
		err = ps.DeletePin(pinbase.Hash("foo"), pinbase.Hash("baz"))
		if err != nil {
			t.Errorf("failed to delete nonexistent pin: %+v", err)
		}

		pin, err := ps.Pin(pinbase.Hash("foo"), pinbase.Hash("baz"))
		if err != nil {
			t.Errorf("failed to try to get the pin: %+v", err)
		}

		if pin != nil {
			t.Errorf("somehow created a pin while deleting it: %+v", pin)
		}
	}
}

func checkBump(t *testing.T, tag string, expect bool, c <-chan struct{}) {
	select {
	case <-c:
		if expect {
			break
		} else {
			t.Errorf("%s: got an unexpected bump", tag)
		}

	case <-time.After(25 * time.Microsecond):
		if expect {
			t.Errorf("%s: did not get the expected bump", tag)
		} else {
			break
		}
	}
}

func TestPinBackendHappyPath(t *testing.T, pb pinbase.PinBackend, ps pinbase.PinService) {
	reqs := pb.PinRequirements()
	if len(reqs) != 0 {
		t.Errorf("requiremens did not start out empty: %+v", reqs)
	}

	err := ps.CreateParty(&pinbase.PartyCreate{
		ID:          pinbase.Hash("foo"),
		Description: "hello",
	})
	if err != nil {
		t.Errorf("failed to create party: %+v", err)
	}

	reqs = pb.PinRequirements()
	if len(reqs) != 0 {
		t.Errorf("requirements not empty: %+v", reqs)
	}

	checkBump(t, "before pins", false, pb.PinProcessorBump())

	err = ps.CreatePin(
		pinbase.Hash("foo"),
		&pinbase.PinCreate{
			ID:         pinbase.Hash("bar"),
			Aliases:    []string{"something"},
			WantPinned: true,
		},
	)
	if err != nil {
		t.Errorf("failed to create pin: %+v", err)
	}

	checkBump(t, "pin created", true, pb.PinProcessorBump())

	reqs = pb.PinRequirements()
	if !reflect.DeepEqual(
		reqs,
		map[pinbase.Hash]bool{
			pinbase.Hash("bar"): true,
		},
	) {
		t.Errorf("got the wrong requirements: %+v", reqs)
	}

	err = ps.UpdatePin(
		pinbase.Hash("foo"),
		pinbase.Hash("bar"),
		&pinbase.PinEdit{
			Aliases:    []string{"something", "something else"},
			WantPinned: true,
		},
	)
	if err != nil {
		t.Errorf("failed to update pin: %+v", err)
	}

	checkBump(t, "pin aliases updated", false, pb.PinProcessorBump())

	reqs = pb.PinRequirements()
	if !reflect.DeepEqual(
		reqs,
		map[pinbase.Hash]bool{
			pinbase.Hash("bar"): true,
		},
	) {
		t.Errorf("pin requirements changed unexpectedly: %+v", reqs)
	}

	err = ps.UpdatePin(
		pinbase.Hash("foo"),
		pinbase.Hash("bar"),
		&pinbase.PinEdit{
			Aliases:    []string{"something"},
			WantPinned: false,
		},
	)
	if err != nil {
		t.Errorf("failed to update poin: %+v", err)
	}

	checkBump(t, "pin want updated", true, pb.PinProcessorBump())

	reqs = pb.PinRequirements()
	if !reflect.DeepEqual(
		reqs,
		map[pinbase.Hash]bool{
			pinbase.Hash("bar"): false,
		},
	) {
		t.Errorf("pin requirements are wrong: %+v", reqs)
	}

	err = ps.CreatePin(
		pinbase.Hash("foo"),
		&pinbase.PinCreate{
			ID:         pinbase.Hash("baz"),
			Aliases:    []string{"doomed"},
			WantPinned: true,
		},
	)
	if err != nil {
		t.Errorf("failed to create pin: %+v", err)
	}

	checkBump(t, "doomed pin created", true, pb.PinProcessorBump())

	reqs = pb.PinRequirements()
	if !reflect.DeepEqual(
		reqs,
		map[pinbase.Hash]bool{
			pinbase.Hash("bar"): false,
			pinbase.Hash("baz"): true,
		},
	) {
		t.Errorf("pin requirements are wrong: %+v", reqs)
	}

	err = ps.DeletePin(pinbase.Hash("foo"), pinbase.Hash("baz"))
	if err != nil {
		t.Errorf("failed to delete pin: %+v", err)
	}

	checkBump(t, "pin deleted", true, pb.PinProcessorBump())

	reqs = pb.PinRequirements()
	if !reflect.DeepEqual(
		reqs,
		map[pinbase.Hash]bool{
			pinbase.Hash("bar"): false,
			pinbase.Hash("baz"): false,
		},
	) {
		t.Errorf("pin requirements are wrong: %+v", reqs)
	}

	err = ps.UpdatePin(
		pinbase.Hash("foo"),
		pinbase.Hash("bar"),
		&pinbase.PinEdit{
			Aliases:    []string{"everything is about to end"},
			WantPinned: true,
		},
	)
	if err != nil {
		t.Errorf("failed to update pin: %+v", err)
	}

	checkBump(t, "pin updated again", true, pb.PinProcessorBump())

	reqs = pb.PinRequirements()
	if !reflect.DeepEqual(
		reqs,
		map[pinbase.Hash]bool{
			pinbase.Hash("bar"): true,
			pinbase.Hash("baz"): false,
		},
	) {
		t.Errorf("pin requirements are wrong: %+v", reqs)
	}

	err = ps.DeleteParty(pinbase.Hash("foo"))
	if err != nil {
		t.Errorf("failed to delete party: %+v", err)
	}

	checkBump(t, "party deleted", true, pb.PinProcessorBump())

	reqs = pb.PinRequirements()
	if !reflect.DeepEqual(
		reqs,
		map[pinbase.Hash]bool{
			pinbase.Hash("bar"): false,
			pinbase.Hash("baz"): false,
		},
	) {
		t.Errorf("pin requirements are wrong: %+v", reqs)
	}
}
