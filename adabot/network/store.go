package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var cacheFile = "./.cache.%s.json"
var cache = make(map[string]*Floorplan)

// LoadStorer as interface is still WIP
type LoadStorer interface {
	Store(id string) error
	Load(id string) (*Floorplan, error)
}

// Store implements part of LoadStorer and will write the entire graph structure with netID
func (data *RawGraph) Store(netID string) error {
	return write(data, netID)
}

// Store implements part of Storer and writes the entire plan structure to memory.
func (data *Floorplan) Store(planID string) error {
	// in-memory cache
	cache[planID] = data
	return nil
}

// Load retrieves all nodes, for the given ID, from the store.
func (data *RawGraph) Load(netID string) error {
	// very simply load cached file
	content, err := ioutil.ReadFile(fmt.Sprintf(cacheFile, netID))
	if err == nil {
		err = json.Unmarshal(content, data)
		if err != nil {
			return err
		}
	}
	return err
}

// Load retrieves the Floorplan structure from the cache if resident.
func (data *Floorplan) Load(planID string) (*Floorplan, error) {
	if val, ok := cache[planID]; ok {
		data = val
		fmt.Printf("Loaded: num polygons: %d\n", len(data.Polygons))
	} else {
		return nil, fmt.Errorf("Store: Floorplan not in cache at ID %s", planID)
	}
	return data, nil
}
func write(data interface{}, id string) error {
	// quite simply write to file
	toCache, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf(cacheFile, id), toCache, 0644)
	return err
}
