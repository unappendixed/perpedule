package core

import (
	"os"
	"testing"
)

func TestParseICSFile(t *testing.T) {
    const TEST_FILE string = "../../test_files/testcalendar.ics"
    file, err := os.Open(TEST_FILE)
    if err != nil {
        t.Fatal("Failed to open test calendar asset")
    }

    cal, err := parseICS(file)
    if err != nil {
        t.Fatal("Failed to parse object")
    }

    cd := initCalendarData(cal)
    expected := 2
    if len(cd.components) != expected {
        t.Errorf("want %d got %d", expected, len(cd.components))
    }
}
