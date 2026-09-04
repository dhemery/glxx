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
	if _, ok := archive.Relationships[id]; ok {
		return fmt.Errorf("Not implemented: describe relationship")
	}
	if _, ok := archive.Repositories[id]; ok {
		return fmt.Errorf("Not implemented: describe repository")
	}
	if _, ok := archive.Sources[id]; ok {
		return fmt.Errorf("Not implemented: describe source")
	}

	return fmt.Errorf("Unknown ID: %s", id)
}

func describeEvent(a *glx.GLXFile, id string, e *glx.Event) {
	printReportHeader(id)

	printReportItem("Title:", e.Title)
	printReportItem("Type:", e.Type)
	printReportItem("Place:", placeName(a, e.PlaceID))
	printReportItem("Date:", e.Date)

	printSectionHeader("Participants")
	for _, p := range e.Participants {
		printParticipation(a, p.Person, p.Role)
	}
	fmt.Println()
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
	name := personName(a, personID)
	printReportItem(label, name)
}

func printReportHeader(title string) {
	fmt.Printf("=== %s ===\n\n", title)
}

func printReportItem(label string, value any) {
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
