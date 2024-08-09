// +build !debug

package core

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
)

const SHARED_EXAMPLE_FILE = "../../test_files/testcalendar.ics"

// HELPERS

// Attempts to open the calendar file at `filepath` and calls Fatalf on error
func openCalendarFile(t *testing.T, filepath string) *CalendarData {
	t.Helper()

	file, err := os.Open(filepath)
	if err != nil {
		t.Fatal("Failed to open test calendar asset")
	}
    defer func(){ _ = file.Close()}()

	cal, err := parseICS(file)
	if err != nil {
		t.Fatal("Failed to parse object")
	}

	return NewCalendarData(cal)

}

func diffText(t *testing.T, first any, second any) string {
	t.Helper()

    s1 := fmt.Sprint(first)
    s2 := fmt.Sprint(second)

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(s1, s2, false)

	return dmp.DiffPrettyText(diffs)

}

// Writes the calendar data to a file in /tmp with the filename formatted as:
// "perpedule_{suffix}_{currentUnixTime}.ics" and returns the file path
func writeTempICSFile(t *testing.T, suffix string, cd *CalendarData) string {

	t.Helper()

	cal := cd.ToICal()

	ts := fmt.Sprint(time.Now().Unix())
	fileid := fmt.Sprintf("perpedule_%s_%s.ics", suffix, ts)
	filepath := path.Join(os.TempDir(), fileid)
	tempfilepath, err := os.Create(filepath)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %s", err.Error())
	}

	cal.SerializeTo(tempfilepath)

	return filepath
}

func cleanupTempICSFile(t *testing.T, filepath string) {
    t.Helper()

    noclean := os.Getenv("NOCLEAN")

    if noclean == "1" || noclean == "true" {
        t.Logf("Not removing temp file %s because NOCLEAN is set to %s\n", filepath, noclean)
        return
    }
    err := os.Remove(filepath)
    if err != nil {
        fmt.Printf("Failed to cleanup temp file %q: %s\n", filepath, err.Error())
    }
}

func printComponents(t *testing.T, cd *CalendarData) {
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

// Attempts to find and return a timeblock resource from `cd`. Calls t.Fatal on
// any errors
func tryGetProject(t *testing.T, cd *CalendarData, uid string) Project {
    t.Helper()

    comp, found := cd.GetByUid(uid)
    if !found {
        t.Fatalf("Failed to lookup project with uid %q", uid)
    }

    p, success := cd.GetAsProject(comp)
    if !success {
        t.Fatalf("Failed to convert component with uid %q to project", comp)
    }

    return p
}

// Attempts to find and return a timeblock resource from `cd`. Calls t.Fatal on
// any errors
func tryGetTimeBlock(t *testing.T, cd *CalendarData, uid string) TimeBlock {
    t.Helper()

    comp, found := cd.GetByUid(uid)
    if !found {
        t.Fatalf("Failed to lookup timeblock with uid %q", uid)
    }

    tb, success := cd.GetAsTimeBlock(comp)
    if !success {
        t.Fatalf("Failed to convert component with uid %q to timeblock", comp)
    }

    return tb
}

func TestWriteICSFile(t *testing.T) {
	TEST_FILE := SHARED_EXAMPLE_FILE
	cd := openCalendarFile(t, TEST_FILE)
	tempfilepath := writeTempICSFile(t, "testwriteicsfile", cd)
	cd2 := openCalendarFile(t, tempfilepath)
    cleanupTempICSFile(t, tempfilepath)
	if !cd.Equal(cd2) {
		t.Fatalf("original file and new file are not equal: %s",
			diffText(t, cd, cd2),
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


func TestAddProject(t *testing.T) {
    TEST_FILE := SHARED_EXAMPLE_FILE
    cd := openCalendarFile(t, TEST_FILE)
    p, err := NewProject("Test project")
    if err != nil {
        t.Fatalf("failed to create project: %s", err.Error())
    }

    cd.AddProject(p)

    tempfilepath := writeTempICSFile(t, "testaddproject", cd)

    cd2 := openCalendarFile(t, tempfilepath)
    cleanupTempICSFile(t, tempfilepath)

    component, found := cd2.GetByUid(p.Id())
    if !found {
        t.Fatalf("Newly created project did not survive serialization.")
    }

    p2, success := cd2.GetAsProject(component)
    if !success {
        t.Fatalf("Newly created component could not be cast to project.")
    }

    // TODO: Add project.Equal method for more complete test comparisons
    if p.Id() != p2.Id() || p.Name() != p2.Name() {
        t.Fatalf("Project serialization failed: diff %q", diffText(t, p, p2))
    }

}

func TestAddTimeblock(t *testing.T) {
	TEST_FILE := SHARED_EXAMPLE_FILE
	cd := openCalendarFile(t, TEST_FILE)

	nb, err := NewTimeBlock("Test block")
	if err != nil {
		t.Fatalf("failed to create timeblock: %s", err.Error())
	}

	blockid := nb.Uid()

	cd.AddTimeBlock(nb)

	newfile := writeTempICSFile(t, "testaddtimeblock", cd)

	cd = openCalendarFile(t, newfile)
    cleanupTempICSFile(t, newfile)
	component, found := cd.GetByUid(blockid)
	if !found {
		t.Fatalf("Newly created timeblock did not survive serialization.")
	}

	tb, success := cd.GetAsTimeBlock(component)
	if !success {
		t.Fatalf("Newly created component could not be cast to timeblock.")
	}

	if nb.Uid() != tb.Uid() || nb.Name() != tb.Name() {
		t.Fatalf("Timeblock serialization failed: diff %q", diffText(t, tb, nb))
	}

}

func TestAddTimeblockToProject(t *testing.T) {
    TEST_FILE := SHARED_EXAMPLE_FILE
    cd := openCalendarFile(t, TEST_FILE)

    p, err := NewProject("Test project")
    if err != nil {
        t.Fatalf("Failed to create new project: %s", err.Error())
    }

    tb, err := NewTimeBlock("Test block")
    if err != nil {
        t.Fatalf("Failed to create new timeblock: %s", err.Error())
    }

    tb.SetProperties(SetTimeBlockParent(p))

    cd.AddProject(p)
    cd.AddTimeBlock(tb)

    tempfilepath := writeTempICSFile(t, "testaddtimeblocktoproject", cd)

    cd2 := openCalendarFile(t, tempfilepath)
    cleanupTempICSFile(t, tempfilepath)

    tb2 := tryGetTimeBlock(t, cd2, tb.Uid())

    p2, err := tb2.Parent()
    if err != nil {
        t.Fatalf("Could not get timeblock parent: %s", err.Error())
    }

    if p2.Id() != p.Id() {
        t.Fatalf("Wrong parent after serialization: got %q want %q", p.Id(), p2.Id())
    }

}

