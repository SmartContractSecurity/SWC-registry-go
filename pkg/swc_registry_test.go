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

func TestSWC(t *testing.T) {
    swc101Markdown := "# Title \nInteger Overflow and Underflow\n\n## Relationships\n[CWE-682: Incorrect Calculation](https://cwe.mitre.org/data/definitions/682.html) \n\n## Description \n\nAn overflow/underflow happens when an arithmetic operation reaches the maximum or minimum size of a type. For instance if a number is stored in the uint8 type, it means that the number is stored in a 8 bits unsigned number ranging from 0 to 2^8-1. In computer programming, an integer overflow occurs when an arithmetic operation attempts to create a numeric value that is outside of the range that can be represented with a given number of bits – either larger than the maximum or lower than the minimum representable value.\n\n## Remediation\n\nIt is recommended to use vetted safe math libraries for arithmetic operations consistently throughout the smart contract system.\n\n## References \n- [Ethereum Smart Contract Best Practices - Integer Overflow and Underflow](https://consensys.github.io/smart-contract-best-practices/known_attacks/#integer-overflow-and-underflow)\n"
    swc101Title := "Integer Overflow and Underflow"
    swc101Relationships := "[CWE-682: Incorrect Calculation](https://cwe.mitre.org/data/definitions/682.html)"
    swc101Description := "An overflow/underflow happens when an arithmetic operation reaches the maximum or minimum size of a type. For instance if a number is stored in the uint8 type, it means that the number is stored in a 8 bits unsigned number ranging from 0 to 2^8-1. In computer programming, an integer overflow occurs when an arithmetic operation attempts to create a numeric value that is outside of the range that can be represented with a given number of bits – either larger than the maximum or lower than the minimum representable value."
    swc101Remediation := "It is recommended to use vetted safe math libraries for arithmetic operations consistently throughout the smart contract system."

    // get an SWC entry by ID (without online update)
    swc, err := GetSWC("SWC-101", false)
    if err != nil {
        t.Fatalf("Unexpected error while creating SWC-101: %s", err)
    }

    testMarkdown := swc.GetMarkdown()
    if testMarkdown != swc101Markdown {
        t.Errorf("Encountered invalid Markdown for SWC-101. Expected:\n%s \n...but got:\n%s\n", swc101Markdown, testMarkdown)
    }
    testTitle := swc.GetTitle()
    if testTitle != swc101Title {
        t.Errorf("Encountered invalid description for SWC-101. Expected:\n%s \n...but got:\n%s\n", swc101Title, testTitle)
    }
    testRelationships := swc.GetRelationships()
    if testRelationships != swc101Relationships {
        t.Errorf("Encountered invalid description for SWC-101. Expected:\n%s \n...but got:\n%s\n", swc101Relationships, testRelationships)
    }
    testDescription := swc.GetDescription()
    if testDescription != swc101Description {
        t.Errorf("Encountered invalid description for SWC-101. Expected:\n%s \n...but got:\n%s\n", swc101Description, testDescription)
    }
    testRemediation := swc.GetRemediation()
    if testRemediation != swc101Remediation {
        t.Errorf("Encountered invalid description for SWC-101. Expected:\n%s \n...but got:\n%s\n", swc101Remediation, testRemediation)
    }

    // get an SWC entry by ID (with online update)
    swc, err = GetSWC("SWC-101", true)
    if err != nil {
        t.Fatalf("Unexpected error while creating SWC-101: %s", err)
    }

    testMarkdown = swc.GetMarkdown()
    if testMarkdown != swc101Markdown {
        t.Errorf("Encountered invalid Markdown for SWC-101. Expected:\n%s \n...but got:\n%s\n", swc101Markdown, testMarkdown)
    }
    testTitle = swc.GetTitle()
    if testTitle != swc101Title {
        t.Errorf("Encountered invalid description for SWC-101. Expected:\n%s \n...but got:\n%s\n", swc101Title, testTitle)
    }
    testRelationships = swc.GetRelationships()
    if testRelationships != swc101Relationships {
        t.Errorf("Encountered invalid description for SWC-101. Expected:\n%s \n...but got:\n%s\n", swc101Relationships, testRelationships)
    }
    testDescription = swc.GetDescription()
    if testDescription != swc101Description {
        t.Errorf("Encountered invalid description for SWC-101. Expected:\n%s \n...but got:\n%s\n", swc101Description, testDescription)
    }
    testRemediation = swc.GetRemediation()
    if testRemediation != swc101Remediation {
        t.Errorf("Encountered invalid description for SWC-101. Expected:\n%s \n...but got:\n%s\n", swc101Remediation, testRemediation)
    }

    // try to get an SWC with an invalid ID
    swc, err = GetSWC("invalid", false)
    if err == nil {
        t.Error("Expected error for an invalid SWC ID was not thrown.")
    }

    // try to get an SWC entry (w/ online update) with some HTTP error
    oldGithubURL := DefaultGithubURL
    DefaultGithubURL = "https://invalid"
    swc, err = GetSWC("SWC-101", true)
    DefaultGithubURL = oldGithubURL
    if err == nil {
        t.Errorf("Expected error for SWC with HTTP failure was not thrown.")
    }
}