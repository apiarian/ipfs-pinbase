package test

import (
	"testing"

	"github.com/apiarian/ipfs-pinbase/pinbase"
)

func TestPinServiceHappyPath(t *testing.T, ps pinbase.PinService) {
	err := ps.CreateParty(&pinbase.PartyCreate{
		ID:          pinbase.Hash("foo"),
		Description: "hello world",
	})
	if err != nil {
		t.Errorf("did not create party: %+v", err)
	}

	parties, err := ps.Parties()
	if err != nil {
		t.Errorf("did not get parties; %+v", err)
	}

	t.Logf("parties: %+v", parties)
	t.Error("bump")
	/*
		p, err := ps.Parties()
		if err != nil {
			t.Errorf("did not get parties: %+v", err)
		}

		_ = p
	*/
}
