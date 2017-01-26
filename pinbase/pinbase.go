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
