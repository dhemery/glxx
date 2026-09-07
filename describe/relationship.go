package describe

import (
	"fmt"
	"io"

	"github.com/genealogix/glx/go-glx"
)

type Relationship Entity[glx.Relationship]

func (relationship *Relationship) Describe(w io.Writer) {
	id := relationship.id
	r := relationship.entity
	a := relationship.archive.file

	printReportHeader("Relationship", id)
	fmt.Fprintln(w)

	printReportItem("Type:", r.Type)

	for _, p := range r.Participants {
		person := Archive{a}.Person(p.Person)
		printParticipation(person, p.Role)
	}

	printRelationshipEvent(w, a, "Start", r.StartEvent)
	printRelationshipEvent(w, a, "End", r.EndEvent)

	fmt.Fprintln(w)
}

func printRelationshipEvent(w io.Writer, a *glx.GLXFile, label, id string) {
	if id == "" {
		return
	}

	fmt.Fprintln(w)
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
