package party

import "time"

type Party struct {
	hash        string
	kind        Kind
	updated     time.Time
	description string
	pins        map[string]*pin
}

type Kind int

const (
	// NodeKind indicates a party which is an IPFS node
	NodeKind Kind = iota

	// ActorKind indicates a party which is not an IPFS node, usually a
	// human user or some sort of script or tool
	ActorKind

	numKinds
)

func (p Kind) String() string {
	switch p {
	case NodeKind:
		return "node"

	case ActorKind:
		return "actor"

	default:
		return "unknown"
	}
}

func NewParty(hash, description string, kind Kind) *Party {
	return &Party{
		hash:        hash,
		kind:        kind,
		updated:     time.Now(),
		pins:        make(map[string]*pin),
		description: description,
	}
}

func (p *Party) SetDescription(d string) {
	if d != p.description {
		p.updated = time.Now()
	}

	p.description = d
}

func (p *Party) Description() string {
	return p.description
}

func (p *Party) Updated() time.Time {
	return p.updated
}

func (p *Party) Kind() Kind {
	return p.kind
}

func (p *Party) Hash() string {
	return p.hash
}

type Pin struct {
	Hash   string
	Pinned bool
}

type none struct{}

type pin struct {
	hash    string
	pinned  bool
	aliases map[string]none
	scope   Scope
	created time.Time
}

type Scope int

const (
	// PublicScope indicates a pin that can be shared is publicly
	// ackgnowledged as belonging to the associated party
	PublicScope Scope = iota

	// ProtectedScope indicates a pin that can be shared but is not
	// acknowledged as belonging to the associated party. The pin will be
	// attributed to the node's party. If there is only one party in a pinbase,
	// this "protection" is a bit questionable.
	ProtectedScope

	// SecretScope indicates a pin that must not be shared
	SecretScope

	numScopes
)

func (p Scope) String() string {
	switch p {
	case PublicScope:
		return "public"

	case ProtectedScope:
		return "protected"

	case SecretScope:
		return "secret"

	default:
		return "unknown"
	}
}

/*
func (p *Party) AddPin(pin *Pin) bool {
	var changed bool

	x, exists := p.Pins[pin.Hash]
	if exists {
		changed = x.MergePin(pin)
	} else {
		p.Pins[pin.Hash] = pin
		changed = true
	}

	if changed {
		p.Updated = time.Now()
	}

	return changed
}

func (p *Party) RemovePin(pin *Pin) {
	_, exists := p.Pins[pin.Hash]
	if !exists {
		return
	}

	delete(p.Pins, pin.Hash)
	p.Updated = time.Now()
}

func (p *Party) UpdatePin(

func (p *Party) UnPin(hash string) {
	if pin, exists := p.Pins[hash]; exists {
		pin.Pinned = false
		p.Updated = time.Now()
	}
}

func (p *Party) RePin(hash string) {
	if pin, exists := p.Pins[hash]; exists {
		pin.Pinned = true
		p.Updated = time.Now()
	}
}

func NewPin(hash string, pinned bool, scope PinScope, aliases ...string) *Pin {
	p := &Pin{
		Hash:    hash,
		Pinned:  pinned,
		Scope:   scope,
		Created: time.Now(),
		Aliases: nil,
	}

	p.AddAliases(aliases...)

	return p
}

func (p *Pin) MergePin(o *Pin) bool {
	if p.Hash != o.Hash {
		return false
	}

	x := len(p.Aliases)
	p.AddAliases(o.Aliases...)
	return len(p.Aliases) != x
}

func (p *Pin) AddAliases(aliases ...string) {
	l := make(map[string]none)
	for _, a := range p.Aliases {
		l[a] = none{}
	}

	for _, a := range aliases {
		if _, exists := l[a]; !exists {
			p.Aliases = append(p.Aliases, a)
			l[a] = none{}
		}
	}
}

func (p *Pin) RemoveAliases(aliases ...string) {
	l := make(map[string]none)
	for _, a := range aliases {
		l[a] = none{}
	}

	// This makes the optimistic assumption that we are actually have all of the
	// aliases that we've been asked to remove. If this is not true, we'll end up
	// having an aliases slice backed by a larger array than necessary. NBD
	s := len(p.Aliases) - len(aliases)
	if s < 0 {
		s = 0
	}

	f := make([]string, s)
	for _, a := range p.Aliases {
		if _, exists := l[a]; !exists {
			f = append(f, a)
		}
	}
	p.Aliases = f
}
*/
