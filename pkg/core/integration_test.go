package core

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
	model "github.com/unappendixed/perpedule/pkg/core/model"
)

const SHARED_EXAMPLE_FILE = "../../test_files/testcalendar.ics"

// Attempts to open the calendar file at `filepath` and calls Fatalf on error
func openCalendarFile(t *testing.T, filepath string) *model.CalendarData {
	t.Helper()

	file, err := os.Open(filepath)
	if err != nil {
		t.Fatal("Failed to open test calendar asset")
	}

	cal, err := parseICS(file)
	if err != nil {
		t.Fatal("Failed to parse object")
	}

	return model.NewCalendarData(cal)

}

func diffText(t *testing.T, first string, second string) string {
	t.Helper()

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(first, second, false)

	return dmp.DiffPrettyText(diffs)

}

// Writes the calendar data to a file in /tmp with the filename formatted as:
// "perpedule_{suffix}_{currentUnixTime}.ics" and returns the file path
func writeTempICSFile(t *testing.T, suffix string, cd *model.CalendarData) string {

	t.Helper()

	cal := cd.ToICal()

	ts := fmt.Sprint(time.Now().Unix())
	fileid := fmt.Sprintf("perpedule_%s_%s.ics", suffix, ts)
	filepath := path.Join(os.TempDir(), fileid)
	outfile, err := os.Create(filepath)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %s", err.Error())
	}

	cal.SerializeTo(outfile)

	return filepath
}

func TestWriteICSFile(t *testing.T) {
	TEST_FILE := SHARED_EXAMPLE_FILE
	cd := openCalendarFile(t, TEST_FILE)
	outpath := writeTempICSFile(t, "testwriteicsfile", cd)
	cd2 := openCalendarFile(t, outpath)
	if !cd.Equal(cd2) {
		t.Fatalf("original file and new file are not equal: %s",
			diffText(t, fmt.Sprint(cd), fmt.Sprint(cd2)),
		)
	}
}

func TestParseICSFile(t *testing.T) {
	TEST_FILE := SHARED_EXAMPLE_FILE
	cd := openCalendarFile(t, TEST_FILE)
	expected := 2
	if len(cd.Components) != expected {
		t.Errorf("want %d got %d", expected, len(cd.Components))
	}
}

func printComponents(t *testing.T, cd *model.CalendarData) {
	t.Helper()

	t.Log("Components:")
	for _, v := range cd.Components {
		t.Logf("Component id: %s", v.Id())
		for _, v2 := range v.Properties {
			t.Logf("%v\n", v2)
		}
	}

	t.Log("New Components:")
	for _, v := range cd.NewComponents {
		t.Logf("Component %s", v.Id())
		for _, v2 := range v.Properties {
			t.Logf("%v\n", v2)
		}
	}
}

func TestAddTimeblock(t *testing.T) {
	TEST_FILE := SHARED_EXAMPLE_FILE
	cd := openCalendarFile(t, TEST_FILE)

	nb, err := model.NewTimeBlock("Test block")
	if err != nil {
		t.Fatalf("failed to create timeblock: %s", err.Error())
	}

	blockid := nb.Uid()

	cd.AddTimeBlock(nb)
	printComponents(t, cd)

	newfile := writeTempICSFile(t, "testaddtimeblock", cd)

	cd = openCalendarFile(t, newfile)
	component, found := cd.GetByUid(blockid)
	if !found {
		t.Fatalf("Newly created timeblock did not survive serialization.")
	}

	tb, success := cd.GetAsTimeBlock(component)
	if !success {
		t.Fatalf("Newly created component could not be cast to timeblock.")
	}

	if nb.Uid() != tb.Uid() || nb.Name() != tb.Name() {
		t.Fatalf("Timeblock serialization failed: want %v got %v", nb, tb)
	}

}

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
