package bolt

import (
	"bytes"
	"encoding/gob"
	cerrors "errors"
	"log"
	"time"

	"github.com/apiarian/ipfs-pinbase/pinbase"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

var (
	PartiesBucketKey         = []byte("PARTIES")
	PartyBucketDataKey       = []byte("DATA")
	PartyBucketPinsBucketKey = []byte("PINS")
	PinArchiveBucketKey      = []byte("PIN-ARCHIVE")
)

type Client struct {
	path string
	db   *bolt.DB
	bump chan struct{}
}

func NewClient(path string) *Client {
	return &Client{
		path: path,
		bump: make(chan struct{}),
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

	_, err = tx.CreateBucketIfNotExists(PinArchiveBucketKey)
	if err != nil {
		return errors.Wrap(err, "create pin archive bucket")
	}

	return nil
}

func (c *Client) PinService() pinbase.PinService {
	return &PinService{
		db:   c.db,
		bump: c.bump,
	}
}

func (c *Client) PinBackend() pinbase.PinBackend {
	return &PinService{
		db:   c.db,
		bump: c.bump,
	}
}

type PinService struct {
	db   *bolt.DB
	bump chan struct{}
}

//
// pinbase.PinService implementation
//

func getPartiesBucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	p := tx.Bucket(PartiesBucketKey)
	if p == nil {
		return nil, errors.New("no parties bucket found")
	}

	return p, nil
}

type partyStorage struct {
	Description string
}

func extractPartyStorage(party *bolt.Bucket) (*partyStorage, error) {
	partyData := party.Get(PartyBucketDataKey)
	if partyData == nil {
		return nil, errors.New("did not get party data")
	}

	var p partyStorage
	err := gob.NewDecoder(bytes.NewBuffer(partyData)).Decode(&p)
	if err != nil {
		return nil, errors.Wrap(err, "decode party data")
	}

	return &p, nil
}

func writePartyStorage(party *bolt.Bucket, p *partyStorage) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)

	err := enc.Encode(p)
	if err != nil {
		return errors.Wrap(err, "encode party data")
	}

	err = party.Put(PartyBucketDataKey, b.Bytes())
	if err != nil {
		return errors.Wrap(err, "put party data")
	}

	return nil
}

func getArchiveBucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	a := tx.Bucket(PinArchiveBucketKey)
	if a == nil {
		return nil, errors.New("no archive bucket found")
	}

	return a, nil
}

var sentinel = []byte("x")

func (ps *PinService) Parties() ([]*pinbase.PartyView, error) {
	var list []*pinbase.PartyView

	err := ps.db.View(func(tx *bolt.Tx) error {
		parties, err := getPartiesBucket(tx)
		if err != nil {
			return err
		}

		c := parties.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v != nil {
				return errors.New("found a non-bucket party")
			}

			party := parties.Bucket(k)
			if party == nil {
				return errors.New("did not get party bucket")
			}

			ps, err := extractPartyStorage(party)
			if err != nil {
				return err
			}

			pv := &pinbase.PartyView{
				ID:          pinbase.Hash(k),
				Description: ps.Description,
			}

			list = append(list, pv)
		}
		return nil
	})

	return list, err
}

func (ps *PinService) Party(h pinbase.Hash) (*pinbase.PartyView, error) {
	var p *pinbase.PartyView

	err := ps.db.View(func(tx *bolt.Tx) error {
		parties, err := getPartiesBucket(tx)
		if err != nil {
			return err
		}

		party := parties.Bucket([]byte(h))
		if party == nil {
			// no party is not an error, just a nil party
			return nil
		}

		ps, err := extractPartyStorage(party)
		if err != nil {
			return err
		}

		p = &pinbase.PartyView{
			ID:          h,
			Description: ps.Description,
		}

		return nil
	})

	return p, err
}

func (ps *PinService) CreateParty(p *pinbase.PartyCreate) error {
	return ps.db.Update(func(tx *bolt.Tx) error {
		parties, err := getPartiesBucket(tx)
		if err != nil {
			return err
		}

		partyKey := []byte(p.ID)

		existingParty := parties.Bucket(partyKey)
		if existingParty != nil {
			return errors.New("party already exists")
		}

		newParty, err := parties.CreateBucket(partyKey)
		if err != nil {
			return errors.Wrap(err, "create party bucket")
		}

		err = writePartyStorage(
			newParty,
			&partyStorage{
				Description: p.Description,
			},
		)
		if err != nil {
			return errors.Wrap(err, "put party data")
		}

		_, err = newParty.CreateBucket(PartyBucketPinsBucketKey)
		if err != nil {
			return errors.Wrap(err, "create party-pins bucket")
		}

		return nil
	})
}

func (ps *PinService) DeleteParty(h pinbase.Hash) error {
	var pinsDeleted bool

	err := ps.db.Update(func(tx *bolt.Tx) error {
		parties, err := getPartiesBucket(tx)
		if err != nil {
			return err
		}

		partyKey := []byte(h)

		party := parties.Bucket(partyKey)
		if party == nil {
			// deleting something that does not exist is not an error
			return nil
		}

		pins := party.Bucket(PartyBucketPinsBucketKey)
		if pins == nil {
			return errors.New("did not get a pins bucket")
		}

		oldPins := make(map[pinbase.Hash]struct{})

		c := pins.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			oldPins[pinbase.Hash(k)] = struct{}{}
		}

		err = parties.DeleteBucket(partyKey)
		if err != nil {
			return errors.Wrap(err, "delete party bucket")
		}

		archive, err := getArchiveBucket(tx)
		if err != nil {
			return err
		}

		for h, _ := range oldPins {
			err := archive.Put([]byte(h), sentinel)
			if err != nil {
				return errors.Wrapf(err, "archive pin %s", h)
			}
		}

		pinsDeleted = len(oldPins) > 0

		return nil
	})

	if err != nil {
		return err
	}

	if pinsDeleted {
		go func(c chan<- struct{}) { c <- struct{}{} }(ps.bump)
	}

	return nil
}

func (ps *PinService) UpdateParty(h pinbase.Hash, p *pinbase.PartyEdit) error {
	return ps.db.Update(func(tx *bolt.Tx) error {
		parties, err := getPartiesBucket(tx)
		if err != nil {
			return err
		}

		party := parties.Bucket([]byte(h))
		if party == nil {
			return errors.New("could not find party")
		}

		ps, err := extractPartyStorage(party)
		if err != nil {
			return err
		}

		ps.Description = p.Description

		return writePartyStorage(party, ps)
	})
}

func getPinsBucket(tx *bolt.Tx, h pinbase.Hash) (*bolt.Bucket, error) {
	parties, err := getPartiesBucket(tx)
	if err != nil {
		return nil, err
	}

	party := parties.Bucket([]byte(h))
	if party == nil {
		return nil, errors.New("could not find party")
	}

	pins := party.Bucket(PartyBucketPinsBucketKey)
	if pins == nil {
		return nil, errors.New("did not get a pins bucket")
	}

	return pins, nil
}

type pinStorage struct {
	Aliases          []string
	WantPinned       bool
	Status           pinbase.PinStatus
	LastErrorMessage string
}

func extractPinStorage(data []byte) (*pinStorage, error) {
	var p pinStorage
	err := gob.NewDecoder(bytes.NewBuffer(data)).Decode(&p)
	if err != nil {
		return nil, errors.Wrap(err, "decode pin data")
	}

	return &p, nil
}

func writePinStorage(pins *bolt.Bucket, h pinbase.Hash, p *pinStorage) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)

	err := enc.Encode(p)
	if err != nil {
		return errors.Wrap(err, "encode pin data")
	}

	err = pins.Put([]byte(h), b.Bytes())
	if err != nil {
		return errors.Wrap(err, "put pin data")
	}

	return nil
}

func (ps *PinService) Pins(partyID pinbase.Hash) ([]*pinbase.PinView, error) {
	var list []*pinbase.PinView

	err := ps.db.View(func(tx *bolt.Tx) error {
		pins, err := getPinsBucket(tx, partyID)
		if err != nil {
			return err
		}

		c := pins.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v == nil {
				return errors.New("found a bucket pin")
			}

			ps, err := extractPinStorage(v)
			if err != nil {
				return err
			}

			pv := &pinbase.PinView{
				ID:         pinbase.Hash(k),
				Aliases:    ps.Aliases,
				WantPinned: ps.WantPinned,
				Status:     ps.Status,
				LastError:  nil,
			}

			if ps.LastErrorMessage != "" {
				pv.LastError = cerrors.New(ps.LastErrorMessage)
			}

			list = append(list, pv)
		}

		return nil
	})

	return list, err
}

func (ps *PinService) Pin(partyID, pinID pinbase.Hash) (*pinbase.PinView, error) {
	var p *pinbase.PinView

	err := ps.db.View(func(tx *bolt.Tx) error {
		pins, err := getPinsBucket(tx, partyID)
		if err != nil {
			return err
		}

		pin := pins.Get([]byte(pinID))
		if pin == nil {
			// no pin is not an error, just a nil pin
			return nil
		}

		ps, err := extractPinStorage(pin)
		if err != nil {
			return err
		}

		p = &pinbase.PinView{
			ID:         pinID,
			Aliases:    ps.Aliases,
			WantPinned: ps.WantPinned,
			Status:     ps.Status,
			LastError:  nil,
		}

		if ps.LastErrorMessage != "" {
			p.LastError = cerrors.New(ps.LastErrorMessage)
		}

		return nil
	})

	return p, err
}

func (ps *PinService) CreatePin(partyID pinbase.Hash, pc *pinbase.PinCreate) error {
	err := ps.db.Update(func(tx *bolt.Tx) error {
		pins, err := getPinsBucket(tx, partyID)
		if err != nil {
			return err
		}

		pinKey := []byte(pc.ID)

		existingPin := pins.Get(pinKey)
		if existingPin != nil {
			return errors.New("pin already exists")
		}

		return writePinStorage(
			pins,
			pc.ID,
			&pinStorage{
				Aliases:          pc.Aliases,
				WantPinned:       pc.WantPinned,
				Status:           pinbase.PinPending,
				LastErrorMessage: "",
			},
		)
	})
	if err != nil {
		return err
	}

	go func(c chan<- struct{}) { c <- struct{}{} }(ps.bump)

	return nil
}

func (ps *PinService) DeletePin(partyID, pinID pinbase.Hash) error {
	err := ps.db.Update(func(tx *bolt.Tx) error {
		pins, err := getPinsBucket(tx, partyID)
		if err != nil {
			return err
		}

		pinKey := []byte(pinID)

		pin := pins.Get(pinKey)
		if pin == nil {
			// deleting something that does not exist is not an error
			return nil
		}

		err = pins.Delete(pinKey)
		if err != nil {
			return errors.Wrap(err, "delete pin data")
		}

		archive, err := getArchiveBucket(tx)
		if err != nil {
			return err
		}

		return errors.Wrap(archive.Put(pinKey, sentinel), "archive the pin")
	})

	if err != nil {
		return err
	}

	go func(c chan<- struct{}) { c <- struct{}{} }(ps.bump)

	return nil
}

func (ps *PinService) UpdatePin(partyID, pinID pinbase.Hash, pe *pinbase.PinEdit) error {
	var wantChanged bool

	err := ps.db.Update(func(tx *bolt.Tx) error {
		pins, err := getPinsBucket(tx, partyID)
		if err != nil {
			return err
		}

		pin := pins.Get([]byte(pinID))
		if pin == nil {
			return errors.New("could not find pin")
		}

		ps, err := extractPinStorage(pin)
		if err != nil {
			return err
		}

		if ps.WantPinned != pe.WantPinned {
			wantChanged = true
		}

		ps.Aliases = pe.Aliases
		ps.WantPinned = pe.WantPinned
		ps.Status = pinbase.PinPending
		ps.LastErrorMessage = ""

		return writePinStorage(pins, pinID, ps)
	})

	if err != nil {
		return err
	}

	if wantChanged {
		go func(c chan<- struct{}) { c <- struct{}{} }(ps.bump)
	}

	return nil
}

//
// pinbase.PinBackend implementation
//

func (ps *PinService) PinProcessorBump() <-chan struct{} {
	return ps.bump
}

func (ps *PinService) PinRequirements() map[pinbase.Hash]bool {
	m := make(map[pinbase.Hash]bool)

	err := ps.db.View(func(tx *bolt.Tx) error {
		parties, err := getPartiesBucket(tx)
		if err != nil {
			return err
		}

		archive, err := getArchiveBucket(tx)
		if err != nil {
			return err
		}

		partiesC := parties.Cursor()

		for partyK, partyV := partiesC.First(); partyK != nil; partyK, partyV = partiesC.Next() {
			if partyV != nil {
				log.Printf("non-bucket party found at %s", partyK)
				continue
			}

			party := parties.Bucket(partyK)
			if party == nil {
				log.Printf("did not get bucket for party %s", partyK)
				continue
			}

			pins := party.Bucket(PartyBucketPinsBucketKey)
			if pins == nil {
				log.Printf("did not get pins bucket for party %s", partyK)
				continue
			}

			pinsC := pins.Cursor()

			for pinK, pinV := pinsC.First(); pinK != nil; pinK, pinV = pinsC.Next() {
				pinHash := pinbase.Hash(pinK)

				if m[pinHash] {
					// if the pin is already marked, no need to extract its data
					continue
				}

				ps, err := extractPinStorage(pinV)
				if err != nil {
					log.Printf("failed to extract data for pin %s under party %s", pinK, partyK)
					continue
				}

				m[pinHash] = ps.WantPinned
			}
		}

		archiveC := archive.Cursor()

		for pinK, _ := archiveC.First(); pinK != nil; pinK, _ = archiveC.Next() {
			pinHash := pinbase.Hash(pinK)

			if _, exists := m[pinHash]; !exists {
				m[pinHash] = false
			}
		}

		return nil
	})
	if err != nil {
		log.Printf("error in bolt transaction: %s", err)
	}

	return m
}

func (ps *PinService) NotifyPin(pinID pinbase.Hash, s *pinbase.PinBackendState) {
	pinKey := []byte(pinID)

	err := ps.db.Update(func(tx *bolt.Tx) error {
		var partiesWithPin []pinbase.Hash

		// first go through the parties to find the ones that need updating

		parties, err := getPartiesBucket(tx)
		if err != nil {
			return err
		}

		partiesC := parties.Cursor()

		for partyK, partyV := partiesC.First(); partyK != nil; partyK, partyV = partiesC.Next() {
			if partyV != nil {
				log.Printf("non-bucket party found at %s", partyK)
				continue
			}

			party := parties.Bucket(partyK)
			if party == nil {
				log.Printf("did not get bucket for party %s", partyK)
				continue
			}

			pins := party.Bucket(PartyBucketPinsBucketKey)
			if pins == nil {
				log.Printf("did not get pins bucket for party %s", partyK)
				continue
			}

			pin := pins.Get(pinKey)
			if pin != nil {
				partiesWithPin = append(partiesWithPin, pinbase.Hash(partyK))
			}
		}

		// now update the pins in the parties that need updating

		for _, partyID := range partiesWithPin {
			pins, err := getPinsBucket(tx, partyID)
			if err != nil {
				log.Printf("did not get pins for party %s", partyID)
				continue
			}

			pin := pins.Get(pinKey)
			if pin == nil {
				log.Printf("most suprisingly did not find pin %s for party %s", pinID, partyID)
				continue
			}

			ps, err := extractPinStorage(pin)
			if err != nil {
				log.Printf("failed to extract pin %s for party %s: %s", pinID, partyID, err)
				continue
			}

			ps.Status = s.Status
			if s.LastError == nil {
				ps.LastErrorMessage = ""
			} else {
				ps.LastErrorMessage = s.LastError.Error()
			}

			err = writePinStorage(pins, pinID, ps)
			if err != nil {
				log.Printf("failed to store pin %s for party %s: %s", pinID, partyID, err)
				continue
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("error in bolt transaction: %s", err)
	}
}

var _ pinbase.PinService = &PinService{}
var _ pinbase.PinBackend = &PinService{}
