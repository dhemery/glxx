package describe

import (
	"github.com/genealogix/glx/go-glx"
)

type Event Entity[glx.Event]

func (event *Event) Describe(r Report) {

	r.Item("Title:", event.entity.Title)
	r.Item("Type:", event.entity.Type)
	r.Reference(event.Place())
	r.Item("Date:", event.entity.Date.String())

	r.BeginSection("Participants")
	for _, p := range event.Participants() {
		r.Reference(p)
	}
}

func (event *Event) ID() string {
	return event.id
}

func (event *Event) Name() string {
	panic("unimplemented")
}

func (event *Event) Label() string {
	return "Event"
}

func (e *Event) Place() *Place {
	return e.archive.Place(e.entity.PlaceID)
}

func (e *Event) Participants() []Participant {
	return participants(e.archive, e.entity.Participants)
}
