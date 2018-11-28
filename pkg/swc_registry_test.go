package swc

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateRegistryFromFile(t *testing.T) {
	r := GetRegistry()

	// load SWC data from default path
	err := r.UpdateRegistryFromFile()
	if err != nil {
		t.Fatalf("Loading from file failed. Got error %s", err)
	}

	// check if SWC keys are all there
	keyNumber := len(r.data)
	if keyNumber != 29 {
		t.Errorf("Expected registry to hold 29 SWC IDs but got %d", keyNumber)
	}

	// load from custom path
	err = r.UpdateRegistryFromFile("testdata/swc-definition-small.json")
	if err != nil {
		t.Fatalf("Loading from file failed. Got error %s", err)
	}
	keyNumber = len(r.data)
	if keyNumber != 1 {
		t.Errorf("Expected registry to hold one SWC ID but got %d", keyNumber)
	}

	// try to load from invalid path
	err = r.UpdateRegistryFromFile("doesnt/exist.json")
	if err == nil {
		t.Errorf("Expected error for invalid path not thrown")
	}

	// check that previous SWC keys are still there
	keyNumber = len(r.data)
	if keyNumber != 1 {
		t.Errorf("Expected registry to hold 29 SWC IDs but got %d", keyNumber)
	}

	// try to load invalid JSON file
	err = r.UpdateRegistryFromFile("testdata/swc-definition-invalid.json")
	if err == nil {
		t.Errorf("Expected error for invalid JSON file not thrown")
	}

	// check that previous SWC keys are still there
	keyNumber = len(r.data)
	if keyNumber != 1 {
		t.Errorf("Expected registry to hold 29 SWC IDs but got %d", keyNumber)
	}
}

func TestUpdateRegistryFromURL(t *testing.T) {
    // test server to mock invalid JSON response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"invalid"}`)
	}))
	defer ts.Close()

	r := GetRegistry()

	// load SWC data from default URL
	err := r.UpdateRegistryFromURL()
	if err != nil {
		t.Fatalf("Loading from URL failed. Got error %s", err)
	}

	// check if SWC keys are all there
	keyNumber := len(r.data)
	if keyNumber != 29 {
		t.Errorf("Expected registry to hold 29 SWC IDs but got %d", keyNumber)
	}

	// try to load from invalid JSON at valid URL
    oldDefaultURL := DefaultGithubURL
    DefaultGithubURL = ts.URL
    err = r.UpdateRegistryFromURL()
    if err == nil {
        t.Errorf("Expected error fotr invalid JSON at valid URL not thrown")
    }

    // check that previous SWC keys are still there
	keyNumber = len(r.data)
	if keyNumber != 29 {
		t.Errorf("Expected registry to hold 29 SWC IDs but got %d", keyNumber)
	}
    DefaultGithubURL = oldDefaultURL

	// try to load from invalid URL
	err = r.UpdateRegistryFromURL("https://example.com/invalid.json")
	if err == nil {
		t.Errorf("Expected error for invalid URL not thrown")
	}

	// check that previous SWC keys are still there
	keyNumber = len(r.data)
	if keyNumber != 29 {
		t.Errorf("Expected registry to hold 29 SWC IDs but got %d", keyNumber)
	}
}
