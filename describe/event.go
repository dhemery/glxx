package describe

import (
	"fmt"

	"github.com/genealogix/glx/go-glx"
)

type Event Entity[glx.Event]

func (event *Event) Describe() {
	id := event.id
	e := event.entity
	a := event.archive.file

	printReportHeader("Event", id)
	fmt.Println()

	printReportItem("Title:", e.Title)
	printReportItem("Type:", e.Type)
	printPlaceReference(a, "Place:", e.PlaceID)
	printReportItem("Date:", e.Date.String())

	fmt.Println()
	printSectionHeader("Participants")
	for _, p := range e.Participants {
		person := Archive{a}.Person(p.Person)
		printParticipation(person, p.Role)
	}

	fmt.Println()
}
