package core

import (
	"os"
	"testing"
    model "github.com/unappendixed/perpedule/pkg/core/model"
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

    cd := model.InitCalendarData(cal)
    expected := 2
    if len(cd.Components) != expected {
        t.Errorf("want %d got %d", expected, len(cd.Components))
    }
}
