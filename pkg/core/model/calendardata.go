package core

import (
	"fmt"

	"github.com/arran4/golang-ical"
)

type CalendarData struct {
    Components []ics.ComponentBase
    NewComponents []ics.ComponentBase
    UidLookup map[string]int
}

func NewCalendarData(cal *ics.Calendar) *CalendarData {
    events := cal.Events()
    fmt.Printf("Events: %d", len(events))
    todos := cal.Todos()

    cd := CalendarData{}

    cd.Components = make([]ics.ComponentBase, 0, len(events) + len(todos))
    cd.NewComponents = make([]ics.ComponentBase, 0, 20)
    cd.UidLookup = make(map[string]int, len(events) + len(todos) * 2)

    for _, v := range events {
        idx := len(cd.Components)
        uid := v.Id()
        cd.Components = append(cd.Components, v.ComponentBase)
        cd.UidLookup[uid] = idx
    }

    for _, v := range todos {
        if v == nil {
            continue
        }
        idx := len(cd.Components)
        uid := v.Id()
        cd.Components = append(cd.Components, v.ComponentBase)
        cd.UidLookup[uid] = idx
    }

    return &cd
}

func (cd *CalendarData) GetByUid(uid string) (ics.ComponentBase, bool) {
    idx, found := cd.UidLookup[uid]
    if !found {
        return ics.ComponentBase{}, false
    }

    if idx >= len(cd.Components) {
        idx = idx - len(cd.Components)
        if idx > len(cd.NewComponents) {
            return ics.ComponentBase{}, false
        }

        return cd.NewComponents[idx], true
    }

    return cd.Components[idx], true
}

func (cd *CalendarData) Projects() []Project {
    out := make([]Project, 0, (len(cd.Components) / 5) + (len(cd.NewComponents) / 5))
    for _, v := range cd.Components {
        if p, success := cd.GetAsProject(v); success {
            out = append(out, p)
        }
    }

    for _, v := range cd.NewComponents {
        if p, success := cd.GetAsProject(v); success {
            out = append(out, p)
        }
    }

    return out
}

func (cd *CalendarData) AddProject(p Project) {
    p.calendarData = cd

    cd.NewComponents = append(cd.NewComponents, p.inner)
}

func (cd *CalendarData) GetAsProject(c ics.ComponentBase) (Project, bool) {
	prop := c.GetProperty(ics.ComponentProperty(PPDType))
    out := Project{}

    if prop == nil || prop.Value != PPDTypeProject {
        return out, false
    }

    out.inner = c
    out.calendarData = cd

    return out, true
}

func (cd *CalendarData) TimeBlocks() []TimeBlock {
    out := make([]TimeBlock, 0, (len(cd.Components) / 5) + (len(cd.NewComponents) / 5))
    for _, v := range cd.Components {
        if tb, success := cd.GetAsTimeBlock(v); success {
            out = append(out, tb)
        }
    }

    for _, v := range cd.NewComponents {
        if tb, success := cd.GetAsTimeBlock(v); success {
            out = append(out, tb)
        }
    }

    return out
}

func (cd *CalendarData) AddTimeBlock(tb TimeBlock) {
    cd.NewComponents = append(cd.NewComponents, tb.inner)
}

func (cd *CalendarData) GetAsTimeBlock(c ics.ComponentBase) (TimeBlock, bool) {
    
	prop := c.GetProperty(ics.ComponentProperty(PPDType))
    out := TimeBlock{}

    if prop == nil || prop.Value != PPDTypeProject {
        return out, false
    }

    parentIdProp := c.GetProperty(ics.ComponentProperty(ics.PropertyRelatedTo))
    if parentIdProp != nil {
        out.parentId = parentIdProp.Value
    }

    out.inner = c
    out.calendarData = cd

    return out, true
}
