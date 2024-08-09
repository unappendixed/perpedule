package core

import (
	"github.com/arran4/golang-ical"
	"github.com/google/uuid"
)

type Project struct {
	inner        ics.GeneralComponent
    calendarData *CalendarData
}

type projectPropertyOption func(*Project)

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

	out.inner = ics.GeneralComponent{ComponentBase: todo.ComponentBase, Token: ICSTokenTodo}
	return out, nil
}

func (p Project) Id() string {
    return p.inner.Id()
}

func (p Project) Name() string {
    prop := p.inner.GetProperty(ics.ComponentProperty(ics.PropertyName))
    if prop != nil {
        return prop.Value
    }

    return ""
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

func (p *Project) SetProperties(opts ...projectPropertyOption) {
    for _, v := range opts {
        v(p)
    }
}

func SetProjectName(name string) projectPropertyOption {
    return func(p *Project) {
        p.inner.SetProperty(ics.ComponentProperty(ics.PropertyName), name)
    }
}

func SetProjectDescription(desc string) projectPropertyOption {
    return func(p *Project) {
        p.inner.SetProperty(ics.ComponentProperty(ics.PropertyName), desc)
    }
}
