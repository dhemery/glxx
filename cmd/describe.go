package cmd

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/dhemery/glxx/load"
	"github.com/genealogix/glx/go-glx"
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe an entity",
	Long:  "Describe an entity",
	RunE:  describe,
	Args:  cobra.ExactArgs(1),
}

func describe(_ *cobra.Command, ids []string) error {
	archive, err := load.Load(archivePath)
	if err != nil {
		return err
	}

	id := ids[0]
	if _, ok := archive.Assertions[id]; ok {
		return fmt.Errorf("Not implemented: describe assertions")
	}
	if _, ok := archive.Citations[id]; ok {
		return fmt.Errorf("Not implemented: describe citations")
	}
	if e, ok := archive.Events[id]; ok {
		describeEvent(archive, id, e)
		return nil
	}
	if _, ok := archive.Media[id]; ok {
		return fmt.Errorf("Not implemented: describe media")
	}
	if _, ok := archive.Persons[id]; ok {
		return fmt.Errorf("Not implemented: describe person")
	}
	if _, ok := archive.Places[id]; ok {
		return fmt.Errorf("Not implemented: describe place")
	}
	if r, ok := archive.Relationships[id]; ok {
		describeRelationship(archive, id, r)
		return nil
	}
	if r, ok := archive.Repositories[id]; ok {
		describeRepository(id, r)
		return nil
	}
	if s, ok := archive.Sources[id]; ok {
		describeSource(archive, id, s)
		return nil
	}

	return fmt.Errorf("Unknown ID: %s", id)
}

func describeEvent(a *glx.GLXFile, id string, e *glx.Event) {
	printReportHeader("Event", id)

	printReportItem("Title:", e.Title)
	printReportItem("Type:", e.Type)
	printPlaceReference(a, "Place:", e.PlaceID)
	printReportItem("Date:", e.Date.String())

	printSectionHeader("Participants")
	for _, p := range e.Participants {
		printParticipation(a, p.Person, p.Role)
	}

	fmt.Println()
}

func describeRelationship(a *glx.GLXFile, id string, r *glx.Relationship) {
	printReportHeader("Relationship", id)

	printReportItem("Type:", r.Type)

	for _, p := range r.Participants {
		printParticipation(a, p.Person, p.Role)
	}

	printRelationshipEvent(a, "Start", r.StartEvent)
	printRelationshipEvent(a, "End", r.EndEvent)

	fmt.Println()
}

func describeRepository(id string, r *glx.Repository) {
	printReportHeader("Repository", id)

	printReportItem("Name:", r.Name)
	printReportItem("Type:", r.Type)
	printReportItem("Address:", r.Address)
	printReportItem("City:", r.City)
	printReportItem("State:", r.State)
	printReportItem("Postal Code:", r.PostalCode)
	printReportItem("Country:", r.Country)
	printReportItem("Website:", r.Website)

	fmt.Println()
}

func describeSource(a *glx.GLXFile, id string, s *glx.Source) {
	printReportHeader("Source", id)

	printReportItem("Title:", s.Title)
	for _, author := range s.Authors {
		printReportItem("Author:", author)
	}
	printReportItem("Date:", s.Date.String())
	printReportItem("Language:", s.Language)

	// printReportItem("Media:", s.Title)
	// printReportItem("Repository:", s.Title)

	fmt.Println()
}

func printRelationshipEvent(a *glx.GLXFile, label, id string) {
	if id == "" {
		return
	}

	printSectionHeader(label + " Event: " + id)

	e, ok := a.Events[id]
	if !ok { // Probably can't happen. Validation would have failed.
		printReportItem("Event:", formattedUnknownID(id, "event"))
		return
	}

	printReportItem("Title:", e.Title)
	printReportItem("Type:", e.Type)
	printPlaceReference(a, "Place:", e.PlaceID)
	printReportItem("Date:", e.Date.String())
}

func printPlaceReference(a *glx.GLXFile, label, id string) {
	printReportItem("Place:", placeName(a, id))
	printReportItem("  id:", id)
}

const unspecifiedValue = "—"

func placeName(a *glx.GLXFile, id string) string {
	if id == "" {
		return unspecifiedValue
	}

	p, ok := a.Places[id]
	if !ok {
		return formattedUnknownID(id, "place")
	}

	if p.Name == "" {
		return formattedUnnamedEntity(id, "place")
	}

	if p.ParentID == "" {
		return p.Name
	}

	return p.Name + ", " + placeName(a, p.ParentID)
}

func personName(a *glx.GLXFile, id string) string {
	p, ok := a.Persons[id]
	if !ok {
		return formattedUnknownID(id, "person")
	}

	name := glx.PersonDisplayName(p)
	if name == "" {
		return formattedUnnamedEntity(id, "person")
	}

	return name
}

func printParticipation(a *glx.GLXFile, personID string, role string) {
	label := strings.ToUpper(role[:1]) + role[1:] + ":"
	printPersonReference(a, label, personID)
}

func printPersonReference(a *glx.GLXFile, label, personID string) {
	name := personName(a, personID)
	printReportItem(label, name)
	printReportItem("  id:", personID)
}

func printReportHeader(typ, title string) {
	fmt.Printf("=== %s: %s ===\n\n", typ, title)
}

func printReportItem(label string, value string) {
	if value == "" {
		value = unspecifiedValue
	}
	fmt.Printf("  %-18s%s\n", label, value)
}

func printSectionHeader(title string) {
	const width = 50
	prefix := "\n── " + title + " "
	remaining := max(width-utf8.RuneCountInString(prefix), 2)

	fmt.Println(prefix + strings.Repeat("─", remaining))
}

func formattedUnknownID(id, typ string) string {
	return fmt.Sprintf("unknown %s id %s", typ, id)
}

func formattedUnnamedEntity(id, typ string) string {
	return fmt.Sprintf("unnamed %s id %s", typ, id)
}
