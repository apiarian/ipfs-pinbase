package ipfs

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/ipfs/go-ipfs-api"
	"github.com/pkg/errors"
	"github.com/whyrusleeping/iptb/util"
)

func TestMain(m *testing.M) {
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

	err = iptbutil.IpfsStart(nodes, true, []string{})
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

func addressForNode(n int) (string, error) {
	node, err := iptbutil.LoadNodeN(n)
	if err != nil {
		return "", errors.Wrap(err, "load node")
	}

	addr, err := node.APIAddr()
	return addr, errors.Wrap(err, "get API address")
}

func newShellForNode(n int) (*shell.Shell, error) {
	addr, err := addressForNode(n)
	if err != nil {
		return nil, err
	}

	s := shell.NewShell(addr)
	if !s.IsUp() {
		return nil, errors.New("ipfs node does not seem to be up")
	}

	return s, nil
}
