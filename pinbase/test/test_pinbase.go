package test

import (
	"reflect"
	"testing"

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
		t.Errorf(" failed to try to get the party: %+v", err)
	}

	if party != nil {
		t.Errorf("somehow created a party while deleting it: %+v", party)
	}
}
