package describe

import (
	"fmt"

	"github.com/genealogix/glx/go-glx"
)

type Assertion Entity[glx.Assertion]

func (a *Assertion) Describe() {
	id := a.id
	glxfile := a.archive.file
	e := a.entity
	if e == nil {
		fmt.Println("NO ASSERTION", id)
		return
	}

	printReportHeader("Assertion", id)

	fmt.Println()
	printReportItem("Status:", e.Status)

	fmt.Println()
	printSubjectSection(glxfile, e.Subject)

	fmt.Println()
	printSectionHeader("Conclusion")
	printReportItem("Property:", e.Property)
	if p := e.Participant; p == nil {
		printReportItem("Value:", e.Value)
	} else {
		person := Archive{glxfile}.Person(p.Person)
		printParticipation(person, p.Role)
	}
	printReportItem("Date:", e.Date.String())
	printReportItem("Confidence:", e.Confidence)

	for _, id := range e.Citations {
		fmt.Println()
		printCitationSection(glxfile, id)
	}
	for _, id := range e.Media {
		printMediaReference(glxfile, "Media:", id)
	}
	for _, id := range e.Sources {
		printSourceSection(glxfile, id)
	}

	fmt.Println()
}

func printSubjectSection(a *glx.GLXFile, e glx.EntityRef) {
	switch {
	case e.Event != "":
		printEventSubjectSection(a, e.Event)
	case e.Person != "":
		person := Archive{a}.Person(e.Person)
		printPersonSubjectSection(person)
	case e.Place != "":
		printPlaceSubjectSection(a, e.Place)
	case e.Relationship != "":
		printRelationshipSubjectSection(a, e.Relationship)
	default:
		return
	}
}

func printRelationshipSubjectSection(a *glx.GLXFile, id string) {
	printSectionHeader("Subject Relationship: " + id)
	r, ok := a.Relationships[id]
	if !ok {
		fmt.Println("UNKNOWN RELATIONSHIP")
	}

	printReportItem("Type:", r.Type)

	for _, p := range r.Participants {
		person := Archive{a}.Person(p.Person)
		printParticipation(person, p.Role)
	}

	printRelationshipEvent(a, "Start", r.StartEvent)
	printRelationshipEvent(a, "End", r.EndEvent)
}

func printPersonSubjectSection(p *Person) {
	printSectionHeader("Subject Person: " + p.id)
	printReportItem("Name:", p.Name())
}

func printPlaceSubjectSection(a *glx.GLXFile, id string) {
	printSectionHeader("Subject Place: " + id)
	p := Archive{a}.Place(id)
	printReportItem("Name:", p.Name())
}

func printEventSubjectSection(a *glx.GLXFile, id string) {
	printSectionHeader("Subject Event: " + id)

	e, ok := a.Events[id]
	if !ok {
		fmt.Println("UNKNOWN EVENT")
	}

	printReportItem("Title:", e.Title)
	printReportItem("Type:", e.Type)
	printPlaceReference(a, "Place:", e.PlaceID)
	printReportItem("Date:", e.Date.String())

	for _, p := range e.Participants {
		person := Archive{a}.Person(p.Person)
		printParticipation(person, p.Role)
	}
}
