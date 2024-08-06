package core

import (
	"fmt"
	"github.com/arran4/golang-ical"
    "slices"
)

type CalendarData struct {
    inner *ics.Calendar
    Components []ics.GeneralComponent
    NewComponents []ics.GeneralComponent
    UidLookup map[string]int
}

func (cd *CalendarData) ToICal() ics.Calendar {
    cal := cd.inner
    outComponents := make([]ics.Component, 0, len(cd.Components) + len(cd.NewComponents))
    for _, v := range cal.Components {
        idIndex := slices.IndexFunc[[]ics.IANAProperty](v.UnknownPropertiesIANAProperties(), func(i ics.IANAProperty) bool {
            return i.IANAToken == "UID"
        })

        id := v.UnknownPropertiesIANAProperties()[idIndex].Value

        // Check Components slice
        idInComponents := slices.ContainsFunc[[]ics.GeneralComponent](cd.Components, func(gc ics.GeneralComponent) bool {
            return gc.Id() == id
        })

        if idInComponents{
            continue
        }


        // Check NewComponents slice
        idInNewComponents := slices.ContainsFunc[[]ics.GeneralComponent](cd.NewComponents, func(gc ics.GeneralComponent) bool {
            return gc.Id() == id
        })

        if idInNewComponents {
            continue
        }

        outComponents = append(outComponents, v)

        cal.Components = outComponents

    }

    for _, v := range cd.Components {
        outComponents = append(outComponents, ics.Component(&v))
    }

    for _, v := range cd.NewComponents {
        outComponents = append(outComponents, ics.Component(&v))
    }

    return *cal
}

func NewCalendarData(cal *ics.Calendar) *CalendarData {
    events := cal.Events()
    fmt.Printf("Events: %d", len(events))
    todos := cal.Todos()

    cd := CalendarData{}
    cd.inner = cal

    cd.Components = make([]ics.GeneralComponent, 0, len(events) + len(todos))
    cd.NewComponents = make([]ics.GeneralComponent, 0, 20)
    cd.UidLookup = make(map[string]int, len(events) + len(todos) * 2)

    for _, v := range events {
        idx := len(cd.Components)
        uid := v.Id()
        cd.Components = append(cd.Components, ics.GeneralComponent{ComponentBase: v.ComponentBase, Token: ICSTokenEvent})
        cd.UidLookup[uid] = idx
    }

    for _, v := range todos {
        if v == nil {
            continue
        }
        idx := len(cd.Components)
        uid := v.Id()
        cd.Components = append(cd.Components, ics.GeneralComponent{ComponentBase: v.ComponentBase, Token: ICSTokenTodo})
        cd.UidLookup[uid] = idx
    }

    return &cd
}

func (cd *CalendarData) GetByUid(uid string) (ics.GeneralComponent, bool) {
    idx, found := cd.UidLookup[uid]
    if !found {
        return ics.GeneralComponent{}, false
    }

    if idx >= len(cd.Components) {
        idx = idx - len(cd.Components)
        if idx > len(cd.NewComponents) {
            return ics.GeneralComponent{}, false
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

func (cd *CalendarData) GetAsProject(c ics.GeneralComponent) (Project, bool) {
    if c.Token == ICSTokenTodo &&
    c.GetProperty(ics.ComponentProperty(PPDType)) != nil &&
    c.GetProperty(ics.ComponentProperty(PPDType)).Value == PPDTypeProject {
        return Project{
            inner: c,
            calendarData: cd,
        }, true
    }

    return Project{}, false
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

func (cd *CalendarData) GetAsTimeBlock(c ics.GeneralComponent) (TimeBlock, bool) {
    
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
