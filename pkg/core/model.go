package core

import (
	"fmt"

	ics "github.com/arran4/golang-ical"
)

type calendarData struct {
    components []ics.Component
    newComponents []ics.Component
    uidLookup map[string]int
}

func initCalendarData(cal *ics.Calendar) *calendarData {
    events := cal.Events()
    fmt.Printf("Events: %d", len(events))
    todos := cal.Todos()

    cd := calendarData{}

    cd.components = make([]ics.Component, 0, len(events) + len(todos))
    cd.newComponents = make([]ics.Component, 0, 20)
    cd.uidLookup = make(map[string]int, len(events) + len(todos) * 2)

    for _, v := range events {
        idx := len(cd.components)
        uid := v.Id()
        cd.components = append(cd.components, v)
        cd.uidLookup[uid] = idx
    }

    for _, v := range todos {
        if v == nil {
            continue
        }
        idx := len(cd.components)
        uid := v.Id()
        cd.components = append(cd.components, v)
        cd.uidLookup[uid] = idx
    }

    return &cd
}

func (cd *calendarData) getByUid(uid string) (ics.Component, bool) {
    idx, found := cd.uidLookup[uid]
    if !found {
        return nil, false
    }

    if idx >= len(cd.components) {
        idx = idx - len(cd.components)
        if idx > len(cd.newComponents) {
            return nil, false
        }

        return cd.newComponents[idx], true
    }

    return cd.components[idx], true

}
