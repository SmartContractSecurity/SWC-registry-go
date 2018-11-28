package swc

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
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

// DefaultGithubURL is the default repository URL used for online-loading. It points to the SWC definition JSON.
var DefaultGithubURL = "https://raw.githubusercontent.com/SmartContractSecurity/SWC-registry/master/export/swc-definition.json"

// DefaultFilePath is the default local filesystem path taken when loading from a file.
var DefaultFilePath = "swc-definition.json"

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

// GetSWC returns an SWC entry based on the given SWC ID. If update=true, the most recent SWC file is fetched from the URL at DefaultGithubURL.
func GetSWC(swcID string, update ...bool) (SWC, error) {
    var updateDB bool
    if len(update) >= 1 {
        updateDB = update[0]
    } else {
        // by default don't update the DB
        updateDB = false
    }
    var swc SWC
    r := GetRegistry()
    if updateDB {
        err := r.UpdateRegistryFromURL()
        if err != nil {
            return swc, err
        }
    }
    swc, present := r.data[swcID]
    if !present {
        return swc, errors.New("No SWC found matching the given ID")
    }

    return swc, nil
}

// GetTitle returns the SWC title.
func (s *SWC) GetTitle() string {
    return s.Description.Content.Title
}

// GetRemediation returns the SWC's remediation text.
func (s *SWC) GetRemediation() string {
    return s.Description.Content.Remediation
}

// GetDescription returns the SWC's full-text description.
func (s *SWC) GetDescription() string {
    return s.Description.Content.Description
}

// GetRelationships returns the SWC's equivalent CWE entry formatted as a Markdown link.
func (s *SWC) GetRelationships() string {
    return s.Description.Content.Relationships
}

// GetMarkdown returns a summary of the SWC in Markdown formatting.
func (s *SWC) GetMarkdown() string {
    return s.Description.Markdown
}

func (r *Registry) parseAndUpdate(inputBytes []byte) error {
	var parsedRegistry map[string]Description
	json.Unmarshal(inputBytes, &parsedRegistry)
	if len(parsedRegistry) == 0 {
		return errors.New("Error reading JSON file - invalid or empty")
	}

	// clear current registry and add new data
	r.data = make(map[string]SWC)
	for SWCId, SWCDescription := range parsedRegistry {
		r.data[SWCId] = SWC{ID: SWCId, Description: SWCDescription}
	}
	return nil
}

// UpdateRegistryFromFile loads the SWC definition data from the JSON file in the package's directory.
func (r *Registry) UpdateRegistryFromFile(paths ...string) error {
	var filePath string
	if len(paths) >= 1 {
		filePath = paths[0]
	} else {
		// use local JSON file as default
		filePath = DefaultFilePath
	}
	jsonBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = r.parseAndUpdate(jsonBytes)
	if err != nil {
		return err
	}
	return nil
}

// UpdateRegistryFromURL accesses a JSON file at a remote URL and tries to update the registry with it.
func (r *Registry) UpdateRegistryFromURL(urls ...string) error {
	var url string
	if len(urls) >= 1 {
		url = urls[0]
	} else {
		// use the default Github repo URL
		url = DefaultGithubURL
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	jsonBytes, _ := ioutil.ReadAll(resp.Body)
	err = r.parseAndUpdate(jsonBytes)
	if err != nil {
		return err
	}
	return nil
}
