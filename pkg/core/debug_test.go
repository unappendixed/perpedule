// +build debug

package core

import (
	"os"
	"testing"
)


func TestUnknownComponents(t *testing.T) {
	const TEST_FILE string = "../../test_files/testcalendar.ics"
	file, err := os.Open(TEST_FILE)
	if err != nil {
		t.Fatal("Failed to open test calendar asset")
	}

	cal, err := parseICS(file)
	if err != nil {
		t.Fatal("Failed to parse object")
	}

	for _, v := range cal.Components {
		for _, v2 := range v.UnknownPropertiesIANAProperties() {
			t.Logf("%s\n", v2)
		}
	}
}
