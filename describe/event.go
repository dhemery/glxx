package describe

import (
	"fmt"

	"github.com/genealogix/glx/go-glx"
)

type Event Entity[glx.Event]

func (e *Event) Describe(r Report) {
	r.Begin("Event", e.id)

	r.Item("Title", e.entity.Title)
	r.Item("Type", e.entity.Type)
	e.Place().DescribeAsReference(r, "Place")
	r.Item("Date", e.entity.Date.String())

	r.BeginSection("Participants")
	for _, p := range e.Participants() {
		p.DescribeAsReference(r)
	}

	r.Notes(e.entity.Notes)

	r.End()
}

func (e *Event) DescribeAsSection(r Report, label string) {
	if e == nil {
		return
	}
	r.BeginSection(fmt.Sprintf("%s: %s", label, e.id))

	r.Item("Title", e.entity.Title)
	r.Item("Type", e.entity.Type)
	e.Place().DescribeAsReference(r, "Place")
	r.Item("Date", e.entity.Date.String())

	for _, p := range e.Participants() {
		p.DescribeAsReference(r)
	}
}

func (e *Event) Place() *Place {
	return e.archive.Place(e.entity.PlaceID)
}

func (e *Event) Participants() []*Participant {
	return e.archive.Participants(e.entity.Participants)
}
