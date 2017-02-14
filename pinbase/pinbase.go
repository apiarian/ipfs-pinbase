package pinbase

type Hash string

type PartyCreate struct {
	ID          Hash
	Description string
}

type PartyEdit struct {
	Description string
}

type PartyView struct {
	ID          Hash
	Description string
}

type PinCreate struct {
	ID         Hash
	Aliases    []string
	WantPinned bool
}

type PinEdit struct {
	Aliases    []string
	WantPinned bool
}

type PinView struct {
	ID         Hash
	Aliases    []string
	WantPinned bool
	Status     PinStatus
	LastError  error
}

type PinStatus int

const (
	PinPending PinStatus = iota
	PinPinned
	PinUnpinned
	PinError
	PinFatal
	numPinStatuses
)

type PinService interface {
	Parties() ([]*PartyView, error)
	Party(Hash) (*PartyView, error)

	CreateParty(*PartyCreate) error
	DeleteParty(Hash) error
	UpdateParty(*PartyEdit) error

	Pins(partyID Hash) ([]*PinView, error)
	Pin(partyID, pinID Hash) (*PinView, error)

	CreatePin(partyID Hash, pc *PinCreate) error
	DeletePin(partyID, pinID Hash) error
	UpdatePin(partyID Hash, pe *PinEdit) error
}
