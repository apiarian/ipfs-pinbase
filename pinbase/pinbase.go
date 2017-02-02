package pinbase

type Node struct {
	Hash        string
	Description string
	APIURL      string
}

type NodeStorage interface {
	Nodes() ([]*Node, error)
	NodeForHash(hash string) (*Node, error)
	SaveNode(*Node) error
}

type Party struct {
	Hash        string
	Description string
	Pins        []*Pin
}

type Pin struct {
	Hash    string
	Aliases []string
	Pinned  bool
}

type PartyStorage interface {
	Parties() ([]*Party, error)
	PartyForHash(hash string) (*Party, error)
	SaveParty(*Party) error
}

type Pinner interface {
	Pins() (map[string]struct{}, error)
	Pin(string) error
	Unpin(string) error
}

type Intention struct {
	Party   string
	Object  string
	WantPin bool
}

type InterestTracker interface {
	BootstrapInterest([]Intention)
	UpdateInterest(Intention)
	InterestDigest() map[string]bool
	NotifyState(map[string]struct{})
}

func permanent(err error) bool {
	type p interface {
		IsPermanent() bool
	}

	pErr, ok := err.(p)
	if ok {
		return pErr.IsPermanent()
	}

	return true
}

func ManagePins(
	done <-chan struct{},
	pnr Pinner,
	trkr InterestTracker,
	intentions <-chan Intention,
) {
	type retry struct {
		i Intention
		c int
	}

	var retries []retry

	for {
		select {
		case i := <-intentions:
			trkr.UpdateInterest(i)

			p, err := pnr.Pins()
			if err != nil {
				if permanent(err) {
					panic(err)
				}

				retries = append(retries, retry{i: i})
			}

			for hash, want := range trkr.InterestDigest() {
				_, pinned := p[hash]

				if want && !pinned {
					err = pnr.Pin(hash)
					if err != nil && permanent(err) {
						panic(err)
					}
				}

				if !want && pinned {
					err = pnr.Unpin(hash)
					if err != nil && permanent(err) {
						panic(err)
					}
				}
			}

			p, err = pnr.Pins()
			if err != nil && permanent(err) {
				panic(err)
			}

			trkr.NotifyState(p)

		case <-done:
			return
		}
	}
}
