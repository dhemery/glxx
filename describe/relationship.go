package describe

import (
	"fmt"

	"github.com/genealogix/glx/go-glx"
)

type Relationship Entity[glx.Relationship]

func (relationship *Relationship) Describe() {
	id := relationship.id
	r := relationship.entity
	a := relationship.archive.file

	printReportHeader("Relationship", id)
	fmt.Println()

	printReportItem("Type:", r.Type)

	for _, p := range r.Participants {
		person := Archive{a}.Person(p.Person)
		printParticipation(person, p.Role)
	}

	printRelationshipEvent(a, "Start", r.StartEvent)
	printRelationshipEvent(a, "End", r.EndEvent)

	fmt.Println()
}

func printRelationshipEvent(a *glx.GLXFile, label, id string) {
	if id == "" {
		return
	}

	fmt.Println()
	printSectionHeader(label + " Event: " + id)

	e, ok := a.Events[id]
	if !ok { // Probably can't happen. Validation would have failed.
		printReportItem("Event:", unknown(id, "event"))
		return
	}

	printReportItem("Title:", e.Title)
	printReportItem("Type:", e.Type)
	printPlaceReference(a, "Place:", e.PlaceID)
	printReportItem("Date:", e.Date.String())
}
