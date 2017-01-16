package network

import (
	"github.com/boltdb/bolt"
	"github.com/timshannon/bolthold"
)

// Store abstracts the backend datastore
type Store struct {
	*bolthold.Store
}

var boltdbFile = "./netdb"

// OpenStore will open the backend datastore
func OpenStore() (*Store, error) {
	store, err := bolthold.Open(boltdbFile, 0666, nil)
	return &Store{store}, err
}

// Update will store the given slice of Nodes each with key Node.ID
func (s *Store) Update(data []*Node) error {
	err := s.Bolt().Update(func(tx *bolt.Tx) error {
		for i := range data {
			err := s.TxInsert(tx, data[i].ID, data[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// Query retrieves all nodes, for the given ID, from the store.
func (s *Store) Query(result []*Node, netID string) error {

	err := s.Find(&result, bolthold.Where("NetID").Eq(netID))
	return err
}

// Close will close the backend datastore.
func (s *Store) CloseStore() error {
	return s.Close()
}
