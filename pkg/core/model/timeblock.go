package core

import (
    "github.com/arran4/golang-ical"
)

type TimeBlock struct {
    inner ics.Component
    calendarData *CalendarData
    parentId string
}

