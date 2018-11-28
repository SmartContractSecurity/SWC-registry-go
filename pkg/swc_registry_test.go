package swc

import (
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

	// try to load invalid JSON file
	err = r.UpdateRegistryFromFile("testdata/swc-definition-invalid.json")
	if err == nil {
		t.Errorf("Expected error for invalid JSON file not thrown")
	}
}
