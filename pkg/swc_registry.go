package swc

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sync"
)

// SWC holds the ID and the SWC description.
type SWC struct {
	ID          string
	Description Description
}

// Content holds the detailed SWC data such as title and remediation.
type Content struct {
	Title         string `json:"Title"`
	Relationships string `json:"Relationships"`
	Description   string `json:"Description"`
	Remediation   string `json:"Remediation"`
}

// Description holds the Markdown-formatted description and the content it is comprised of.
type Description struct {
	Markdown string  `json:"markdown"`
	Content  Content `json:"content"`
}

// Registry is a map from SWC ID to the SWC data for faster lookups.
type Registry struct {
	data map[string]SWC
}

var registryInstance *Registry
var once sync.Once

// GetRegistry implements a singleton pattern and fetches the internal Registry instance.
func GetRegistry() *Registry {
	once.Do(func() {
		m := make(map[string]SWC)
		registryInstance = &Registry{data: m}
	})
	return registryInstance
}

// UpdateRegistryFromFile loads the SWC definition data from the JSON file in the package's directory.
func (r *Registry) UpdateRegistryFromFile(path ...string) error {
	var filePath string
	if len(path) >= 1 {
		filePath = path[0]
	} else {
		// use local JSON file as default
		filePath = "swc-definition.json"
	}
	jsonBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var parsedRegistry map[string]Description
	json.Unmarshal(jsonBytes, &parsedRegistry)
	if len(parsedRegistry) == 0 {
		return errors.New("Error reading JSON file - invalid or empty")
	}

	// clear current registry
	r.data = make(map[string]SWC)
	for SWCId, SWCDescription := range parsedRegistry {
		r.data[SWCId] = SWC{ID: SWCId, Description: SWCDescription}
	}
	return nil
}
