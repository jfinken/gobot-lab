package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var cacheFile = "./.cache.%s.json"

// StoreGraph will store the entire graph structure with netID
func StoreGraph(data *RawGraph, netID string) error {
	// quite simply write to file
	toCache, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf(cacheFile, netID), toCache, 0644)
	return err
}

// LoadGraph retrieves all nodes, for the given ID, from the store.
func LoadGraph(data *RawGraph, netID string) error {
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
