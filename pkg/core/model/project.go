package core

import (
	"github.com/arran4/golang-ical"
	"github.com/google/uuid"
)

type Project struct {
	inner        ics.ComponentBase
    calendarData *CalendarData
}

func NewProject(name string) (Project, error) {

	out := Project{
		calendarData: nil,
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return out, err
	}

	todo := ics.NewTodo(id.String())
	todo.AddProperty(ics.ComponentProperty(ics.PropertyName), name)
	todo.AddProperty(ics.ComponentProperty(PPDType), PPDTypeProject)

	out.inner = todo.ComponentBase
	return out, nil
}

func (p Project) Id() string {
    return p.inner.Id()
}

func (p Project) Children() []TimeBlock {

    out := make([]TimeBlock, 0, 10)

    for _, v := range p.calendarData.Components {
        if GetRelatedToId(v) == p.Id() {
            tb, success := p.calendarData.GetAsTimeBlock(v)
            if !success {
                continue
            }

            out = append(out, tb)
        }
    }

    return out
    
}
