package pinbase

import (
	"time"
)

type NodeManager interface {
	AddNode(Node) error
	DeleteNode(hash string) error
	Node(hash string) (Node, error)
	Nodes() ([]Node, error)
}

type Node interface {
	Description() string
	SetDescription(string) error

	Ping() error
	PleasePin(hash, party string)
	PleaseUnpin(hash, party string)
	PinInfo(hash string) *PinInfo
}

type PinInfo struct {
	Hash      string
	Timestamp time.Time
	Status    PinStatus
	Error     error
}

type PinStatus int

const (
	PinPending PinStatus = iota
	PinPinned
	PinUnpinned
	PinTrouble
	PinFailed
	numPinStatuses
)

type Party interface {
	Hash() string
	Description() string
	SetDescription(string) error
}
