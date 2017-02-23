package bolt

import (
	"time"

	"github.com/apiarian/ipfs-pinbase/pinbase"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

var (
	PartiesBucketKey = []byte("PARTIES")
)

type Client struct {
	path string
	db   *bolt.DB
}

func NewClient(path string) *Client {
	return &Client{
		path: path,
	}
}

func (c *Client) Open() error {
	db, err := bolt.Open(
		c.path,
		0600,
		&bolt.Options{
			Timeout: 1 * time.Second,
		},
	)
	if err != nil {
		return errors.Wrapf(err, "connect to db at %s", c.path)
	}

	c.db = db

	err = c.db.Update(setupSchema)
	if err != nil {
		return errors.Wrap(err, "setup the schema")
	}

	return nil
}

func (c *Client) Close() error {
	if c.db != nil {
		return errors.Wrap(c.db.Close(), "close db")
	}

	return nil
}

func setupSchema(tx *bolt.Tx) error {
	_, err := tx.CreateBucketIfNotExists(PartiesBucketKey)
	if err != nil {
		return errors.Wrap(err, "create parties bucket")
	}

	return nil
}

func (c *Client) PinService() pinbase.PinService {
	return &PinService{
		db: c.db,
	}
}

type PinService struct {
	db *bolt.DB
}

//
// pinbase.PinService implementation
//

func (ps *PinService) Parties() ([]*pinbase.PartyView, error) {
	return nil, errors.New("not implemented")
}

func (ps *PinService) Party(h pinbase.Hash) (*pinbase.PartyView, error) {
	return nil, errors.New("not implemented")
}

func (ps *PinService) CreateParty(p *pinbase.PartyCreate) error {
	return errors.New("not implemented")
}

func (ps *PinService) DeleteParty(h pinbase.Hash) error {
	return errors.New("not implemented")
}

func (ps *PinService) UpdateParty(p *pinbase.PartyEdit) error {
	return errors.New("not implemented")
}

func (ps *PinService) Pins(partyID pinbase.Hash) ([]*pinbase.PinView, error) {
	return nil, errors.New("not implemented")
}

func (ps *PinService) Pin(partyID, pinID pinbase.Hash) (*pinbase.PinView, error) {
	return nil, errors.New("not implemented")
}

func (ps *PinService) CreatePin(partyID pinbase.Hash, pc *pinbase.PinCreate) error {
	return errors.New("not implemented")
}

func (ps *PinService) DeletePin(partyID, pinID pinbase.Hash) error {
	return errors.New("not implemented")
}

func (ps *PinService) UpdatePin(partyID pinbase.Hash, pe *pinbase.PinEdit) error {
	return errors.New("not implemented")
}

//
// pinbase.PinBackend implementation
//

func (ps *PinService) PinProcessorBump() <-chan struct{} {
	return make(chan struct{})
}

func (ps *PinService) PinRequirements() map[pinbase.Hash]bool {
	return make(map[pinbase.Hash]bool)
}

func (ps *PinService) NotifyPin(pinID pinbase.Hash, s *pinbase.PinBackendState) {
}

var _ pinbase.PinService = &PinService{}
var _ pinbase.PinBackend = &PinService{}
