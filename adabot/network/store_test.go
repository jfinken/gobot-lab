package network

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/timshannon/bolthold"
)

type Item struct {
	ID       int
	Category string `boltholdIndex:"Category"`
	Created  time.Time
}

// setup creates a slice of Item, and creates a temp File
// in which to store the BoltDB.
func setup() (*os.File, []Item) {
	data := []Item{
		Item{
			ID:       0,
			Category: "blue",
			Created:  time.Now().Add(-4 * time.Hour),
		},
		Item{
			ID:       1,
			Category: "red",
			Created:  time.Now().Add(-3 * time.Hour),
		},
		Item{
			ID:       2,
			Category: "blue",
			Created:  time.Now().Add(-2 * time.Hour),
		},
		Item{
			ID:       3,
			Category: "blue",
			Created:  time.Now().Add(-20 * time.Minute),
		},
	}

	tmpfile, err := ioutil.TempFile(".", "boltholdb")
	if err != nil {
		log.Fatal(err)
	}
	return tmpfile, data
}
func update(s *bolthold.Store, data []Item) error {

	// insert the data in one transaction
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
func TestStoreUpdate(t *testing.T) {

	f, data := setup()

	store, err := bolthold.Open(f.Name(), 0666, nil)
	defer store.Close()
	defer os.Remove(f.Name())

	if err != nil {
		t.Fatalf(err.Error())
	}

	err = update(store, data)

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestStoreFind(t *testing.T) {

	f, data := setup()

	store, err := bolthold.Open(f.Name(), 0666, nil)
	defer store.Close()
	defer os.Remove(f.Name())
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = update(store, data)

	if err != nil {
		t.Fatalf(err.Error())
	}
	// Query: find all items in the blue category that have been created in
	// the past hour
	var result []Item

	err = store.Find(&result, bolthold.Where("Category").Eq("blue").And("Created").Ge(time.Now().Add(-1*time.Hour)))

	if err != nil {
		// handle error
		log.Fatal(err)
	}

	if result[0].ID != 3 {
		t.Errorf("Expected: 3, Got: %d\n", result[0].ID)
	}
}
