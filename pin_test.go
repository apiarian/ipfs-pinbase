package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"

	shell "github.com/apiarian/go-ipfs-api"
	"github.com/pkg/errors"
	"github.com/whyrusleeping/iptb/util"
)

func fatalIfErr(t *testing.T, msg string, err error) {
	if err != nil {
		t.Fatalf("%s: %+v\n", msg, err)
	}
}

func TestMain(m *testing.M) {
	flag.Parse()

	ipfsDir, err := ioutil.TempDir("", "pinbase-test-state-iptb-root")
	if err != nil {
		log.Fatal("failed to create temporary ipfs directory:", err)
	}
	log.Println("temporary ipfs directory:", ipfsDir)

	err = os.Setenv("IPTB_ROOT", ipfsDir)
	if err != nil {
		log.Fatal("failed to set IPTB_ROOT to temproary ipfsDir:", err)
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	ps := 15000 + (rnd.Int()%500)*10
	log.Println("iptb port start:", ps)

	cfg := &iptbutil.InitCfg{
		Count:     2,
		Force:     true,
		Bootstrap: "star",
		PortStart: ps,
		Mdns:      false,
		Utp:       false,
		Override:  "",
		NodeType:  "",
	}
	err = iptbutil.IpfsInit(cfg)
	if err != nil {
		log.Fatal("failed to initialize iptb:", err)
	}

	nodes, err := iptbutil.LoadNodes()
	if err != nil {
		log.Fatal("failed to load nodes:", err)
	}
	defer iptbutil.IpfsKillAll(nodes)

	err = iptbutil.IpfsStart(nodes, true)
	if err != nil {
		for i, n := range nodes {
			killerr := n.Kill()
			if killerr != nil {
				log.Println("failed to kill node", i, ":", killerr)
			} else {
				log.Println("killed node", i)
			}
		}
		log.Fatal("failed to start nodes:", err)
	}

	r := m.Run()

	err = iptbutil.IpfsKillAll(nodes)
	if err != nil {
		log.Print("error killing nodes:", err)
	}

	os.RemoveAll(ipfsDir)

	os.Exit(r)
}

func newShellForNode(n int) (*shell.Shell, error) {
	node, err := iptbutil.LoadNodeN(n)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load node")
	}

	addr, err := node.APIAddr()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get node API address")
	}

	s := shell.NewShell(addr)
	if !s.IsUp() {
		return nil, errors.New("ipfs node does not seem to be up")
	}

	return s, nil
}

func TestOverallPinInterests(t *testing.T) {
	cases := []struct {
		tag    string
		input  map[string]map[string]bool
		output map[string]bool
	}{
		{
			tag: "the basics",
			input: map[string]map[string]bool{
				"objA": map[string]bool{
					"partA": true,
					"partB": false,
				},
				"objB": map[string]bool{
					"partyC": false,
					"partyD": false,
				},
			},
			output: map[string]bool{
				"objA": true,
				"objB": false,
			},
		},
		{
			tag: "party hash is object hash",
			input: map[string]map[string]bool{
				"hashA": map[string]bool{
					"hashA": true,
					"hashB": false,
				},
				"hashB": map[string]bool{
					"hashA": false,
					"hashB": false,
				},
			},
			output: map[string]bool{
				"hashA": true,
				"hashB": false,
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.tag, func(t *testing.T) {
			t.Parallel()
			got := overallPinInterests(c.input)

			if !reflect.DeepEqual(got, c.output) {
				t.Errorf(
					"output (%+v) does not match expected (%+v) for input (%+v)",
					got,
					c.output,
					c.input,
				)
			}
		})
	}
}

func TestPinReconciliation(t *testing.T) {
	s, err := newShellForNode(0)
	fatalIfErr(t, "failed to get a shell for node 0", err)

	hObjX, err := s.Add(bytes.NewBufferString("object X reconciliation"))
	fatalIfErr(t, "failed to add object X", err)

	hObjA, err := s.Add(bytes.NewBufferString("object A reconciliation"))
	fatalIfErr(t, "failed to add object A", err)

	hObjB, err := s.Add(bytes.NewBufferString("object B reconciliation"))
	fatalIfErr(t, "failed to add object B", err)

	cur, err := s.Pins()
	fatalIfErr(t, "failed to get intiial pins", err)

	if _, pinned := cur[hObjA]; !pinned {
		t.Error("object A does not seem to be pinned")
	}

	if _, pinned := cur[hObjB]; !pinned {
		t.Error("object B does not seem to be pinned")
	}

	if _, pinned := cur[hObjX]; !pinned {
		t.Error("object X does not seem to be pinned")
	}

	cases := []struct {
		tag        string
		pins       map[string]bool
		objApinned bool
		objBpinned bool
	}{
		{
			tag: "pin both",
			pins: map[string]bool{
				hObjA: true,
				hObjB: true,
			},
			objApinned: true,
			objBpinned: true,
		},
		{
			tag: "unpin one",
			pins: map[string]bool{
				hObjA: true,
				hObjB: false,
			},
			objApinned: true,
			objBpinned: false,
		},
		{
			tag: "swap pins",
			pins: map[string]bool{
				hObjA: false,
				hObjB: true,
			},
			objApinned: false,
			objBpinned: true,
		},
		{
			tag: "unpin both",
			pins: map[string]bool{
				hObjA: false,
				hObjB: false,
			},
			objApinned: false,
			objBpinned: false,
		},
		{
			tag: "bring them back again",
			pins: map[string]bool{
				hObjA: true,
				hObjB: true,
			},
			objApinned: true,
			objBpinned: true,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.tag, func(t *testing.T) {
			err := reconcilePinNeeds(s, c.pins)
			fatalIfErr(t, "failed to reconcile pins", err)

			cur, err := s.Pins()
			fatalIfErr(t, "failed to get current pins", err)

			if _, pinned := cur[hObjA]; pinned != c.objApinned {
				t.Errorf(
					"expected object A to be pinned(%t) but got (%t)",
					c.objApinned,
					pinned,
				)
			}

			if _, pinned := cur[hObjB]; pinned != c.objBpinned {
				t.Errorf(
					"expected object B to be pinned(%t) but got (%t)",
					c.objBpinned,
					pinned,
				)
			}

			if _, pinned := cur[hObjX]; !pinned {
				t.Errorf("object X has become unpinned")
			}
		})
	}
}

func TestManagePins(t *testing.T) {
	s, err := newShellForNode(0)
	fatalIfErr(t, "failed to get a shell for node 0", err)

	managerS, err := newShellForNode(0)
	fatalIfErr(t, "failed to get a manager shell for node 0", err)

	hObjX, err := s.Add(bytes.NewBufferString("object X management"))
	fatalIfErr(t, "failed to add object X", err)

	hObjA, err := s.AddNoPin(bytes.NewBufferString("object A management"))
	fatalIfErr(t, "failed to add object A", err)

	hObjB, err := s.AddNoPin(bytes.NewBufferString("object B management"))
	fatalIfErr(t, "failed to add object B", err)

	cur, err := s.Pins()
	fatalIfErr(t, "failed to get intiial pins", err)

	if _, pinned := cur[hObjA]; pinned {
		t.Error("object A started out pinned")
	}

	if _, pinned := cur[hObjB]; pinned {
		t.Error("object B started out pinned")
	}

	if _, pinned := cur[hObjX]; !pinned {
		t.Error("object X does not seem to be pinned")
	}

	const (
		hPartyA = "partyA"
		hPartyB = "partyB"
	)

	done := make(chan struct{})
	defer close(done)

	intentions := make(chan PinIntention)

	go ManagePins(managerS, nil, intentions, done)

	cases := []struct {
		tag        string
		intention  PinIntention
		objApinned bool
		objBpinned bool
	}{
		{
			tag: "party A cares about object A",
			intention: PinIntention{
				PartyHash:  hPartyA,
				ObjectHash: hObjA,
				Interested: true,
			},
			objApinned: true,
			objBpinned: false,
		},
		{
			tag: "party B cares about object B",
			intention: PinIntention{
				PartyHash:  hPartyB,
				ObjectHash: hObjB,
				Interested: true,
			},
			objApinned: true,
			objBpinned: true,
		},
		{
			tag: "party B also cares about object A",
			intention: PinIntention{
				PartyHash:  hPartyB,
				ObjectHash: hObjA,
				Interested: true,
			},
			objApinned: true,
			objBpinned: true,
		},
		{
			tag: "party A stopped caring about object A",
			intention: PinIntention{
				PartyHash:  hPartyA,
				ObjectHash: hObjA,
				Interested: false,
			},
			objApinned: true, // party B still cares
			objBpinned: true,
		},
		{
			tag: "party B stopped caring about object A",
			intention: PinIntention{
				PartyHash:  hPartyB,
				ObjectHash: hObjA,
				Interested: false,
			},
			objApinned: false,
			objBpinned: true,
		},
		{
			tag: "party B doesn't care about anything anymore",
			intention: PinIntention{
				PartyHash:  hPartyB,
				ObjectHash: hObjB,
				Interested: false,
			},
			objApinned: false,
			objBpinned: false,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.tag, func(t *testing.T) {
			intentions <- c.intention

			// apparently there needs to be a bit of a delay between pinning
			// something and looking to see if it was actually pinned
			time.Sleep(10 * time.Millisecond)

			cur, err := s.Pins()
			fatalIfErr(t, "failed to get current pins", err)

			if _, pinned := cur[hObjA]; pinned != c.objApinned {
				t.Errorf(
					"expected object A to be pinned(%t) but got (%t)",
					c.objApinned,
					pinned,
				)
			}

			if _, pinned := cur[hObjB]; pinned != c.objBpinned {
				t.Errorf(
					"expected object B to be pinned(%t) but got (%t)",
					c.objBpinned,
					pinned,
				)
			}

			if _, pinned := cur[hObjX]; !pinned {
				t.Errorf("object X has become unpinned")
			}
		})
	}
}
