package core

import (
	"github.com/arran4/golang-ical"
	"github.com/google/uuid"
)

type Project struct {
	inner        ics.ComponentBase
}

func NewProject(name string, cd *CalendarData) (Project, error) {

	out := Project{
		calendarData: cd,
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

func (p Project) Children() []TimeBlock {
    if len(p.children) == 0 {
        return []TimeBlock{}
    }

    // components := make([]TimeBlock, len(p.children))
    // for i := range p.children {
    //     if i > len(p.calendarData.Components) {
    //         // components = append(components, )
    //     }
    // }


    return []TimeBlock{}
    
}
