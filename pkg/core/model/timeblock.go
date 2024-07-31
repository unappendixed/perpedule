package core

import (
	"errors"

	"github.com/arran4/golang-ical"
)

var ParentNotFoundErr error = errors.New("Could not find parent in lookup table")
var InvalidParentTypeErr error = errors.New("Parent was not type project")

type TimeBlock struct {
    inner ics.ComponentBase
    calendarData *CalendarData
    parentId string
}

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
