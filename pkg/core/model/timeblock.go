package core

import (
	"errors"
	"time"

	"github.com/arran4/golang-ical"
)

var ParentNotFoundErr error = errors.New("Could not find parent in lookup table")
var InvalidParentTypeErr error = errors.New("Parent was not type project")

type TimeBlock struct {
    inner ics.GeneralComponent
    calendarData *CalendarData
    parentId string
}

// See SetName, SetDescription, SetParent, SetStart, SetEnd, etc...
type timeBlockPropertyOption func(*TimeBlock)

func (tb TimeBlock) Parent() (*Project, error) {
    c, found := tb.calendarData.GetByUid(tb.parentId)
    if !found {
        return nil, ParentNotFoundErr
    }

    p, success := tb.calendarData.GetAsProject(c)

    if !success {
        return nil, InvalidParentTypeErr
    }

    return &p, nil
}

func (tb *TimeBlock) SetProperties(opts ...timeBlockPropertyOption) {
    for _,v := range opts {
        v(tb)
    }
}

func SetTimeBlockName(name string) timeBlockPropertyOption {
    return func(tb *TimeBlock) {
        tb.inner.SetProperty(ics.ComponentProperty(ics.PropertyName), name)
    }
}

func SetTimeBlockDescription(desc string) timeBlockPropertyOption {
    return func(tb *TimeBlock) {
        tb.inner.SetProperty(ics.ComponentProperty(ics.PropertyName), desc)
    }
}

func SetTimeBlockParent(parentId string) timeBlockPropertyOption {
    return func(tb *TimeBlock) {
        tb.parentId = parentId
    }
}

func SetTimeBlockStart(t time.Time) timeBlockPropertyOption {
    return func(tb *TimeBlock) {
        tb.inner.SetStartAt(t)
    }
}

func SetTimeBlockEnd(t time.Time) timeBlockPropertyOption {
    return func(tb *TimeBlock) {
        tb.inner.SetProperty(ics.ComponentPropertyDtEnd, t.UTC().Format(ICalTimestampFormatUtc))
    }
}
