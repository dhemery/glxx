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
	if a, ok := archive.Assertions[id]; ok {
		describeAssertion(archive, id, a)
		return nil
	}
	if c, ok := archive.Citations[id]; ok {
		describeCitation(archive, id, c)
		return nil
	}
	if e, ok := archive.Events[id]; ok {
		describeEvent(archive, id, e)
		return nil
	}
	if m, ok := archive.Media[id]; ok {
		describeMedia(archive, id, m)
		return nil
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

func describeAssertion(a *glx.GLXFile, id string, e *glx.Assertion) {
	printReportHeader("Assertion", id)

	fmt.Println()
	printReportItem("Status:", e.Status)

	fmt.Println()
	printSectionHeader("Conclusion")

	printSubjectReference(a, "Subject:", e.Subject)
	printReportItem("Property:", e.Property)
	printReportItem("Value:", e.Value)
	printReportItem("Date:", e.Date.String())
	if p := e.Participant; p != nil {
		printParticipation(a, p.Person, p.Role)
	}
	printReportItem("Confidence:", e.Confidence)

	fmt.Println()
	printSectionHeader("Evidence")
	for _, id := range e.Citations {
		printCitationReference(a, "Citation:", id)
	}
	for _, id := range e.Sources {
		printSourceReference(a, "Source:", id)
	}
	for _, id := range e.Media {
		printSourceReference(a, "Media:", id)
	}

	fmt.Println()
}

func describeCitation(a *glx.GLXFile, id string, c *glx.Citation) {
	printReportHeader("Citation", id)
	fmt.Println()

	printSourceReference(a, "Source:", c.SourceID)
	printRepositoryReference(a, "Repository:", c.RepositoryID)
	for _, m := range c.Media {
		printMediaReference(a, "Media:", m)
	}

	fmt.Println()
}

func describeEvent(a *glx.GLXFile, id string, e *glx.Event) {
	printReportHeader("Event", id)
	fmt.Println()

	printReportItem("Title:", e.Title)
	printReportItem("Type:", e.Type)
	printPlaceReference(a, "Place:", e.PlaceID)
	printReportItem("Date:", e.Date.String())

	fmt.Println()
	printSectionHeader("Participants")
	for _, p := range e.Participants {
		printParticipation(a, p.Person, p.Role)
	}

	fmt.Println()
}
func describeMedia(a *glx.GLXFile, id string, m *glx.Media) {
	printReportHeader("Media", id)
	fmt.Println()

	printReportItem("Title:", m.Title)
	printReportItem("URI:", m.URI)
	printReportItem("Type:", m.Type)
	printReportItem("MimeType:", m.MimeType)
	printReportItem("Hash:", m.Hash)
	printReportItem("Date:", m.Date.String())
	printSourceReference(a, "Source:", m.Source)

	fmt.Println()
}

func describeRelationship(a *glx.GLXFile, id string, r *glx.Relationship) {
	printReportHeader("Relationship", id)
	fmt.Println()

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
	fmt.Println()

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
	fmt.Println()

	printReportItem("Title:", s.Title)
	for _, author := range s.Authors {
		printReportItem("Author:", author)
	}
	printReportItem("Date:", s.Date.String())
	printReportItem("Language:", s.Language)

	printRepositoryReference(a, "Repository:", s.RepositoryID)

	for _, m := range s.Media {
		printMediaReference(a, "Media:", m)
	}

	fmt.Println()
}

func printCitationReference(a *glx.GLXFile, label, id string) {
	if id == "" {
		printReportItem(label, unspecifiedValue)
		return
	}

	c, ok := a.Citations[id]
	if !ok {
		printReportItem(label, unknown(id, "citation"))
		return
	}

	printReportLine(label, "")
	printReportItem("  Source:", sourceTitle(a, c.SourceID))
	if c.SourceID != "" {
		printReportItem("    id:", c.SourceID)
	}
	printReportItem("  Repository:", repositoryName(a, c.RepositoryID))
	if c.RepositoryID != "" {
		printReportItem("    id:", c.RepositoryID)
	}
	for _, m := range c.Media {
		printReportItem("  Media:", mediaTitle(a, m))
		printReportItem("    id:", m)
	}
}

func printMediaReference(a *glx.GLXFile, label, id string) {
	printReference(label, id, mediaTitle(a, id))
}

func printPlaceReference(a *glx.GLXFile, label, id string) {
	printReference(label, id, placeName(a, id))
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

func printRepositoryReference(a *glx.GLXFile, label, id string) {
	name := repositoryName(a, id)
	if name == unspecifiedValue {
		printReportItem(label, name)
		return
	}
	printReference(label, id, name)
}

func printSourceReference(a *glx.GLXFile, label, id string) {
	printReference(label, id, sourceTitle(a, id))
}

func printSubjectReference(a *glx.GLXFile, label string, e glx.EntityRef) {
	printReportLine("Subject:", "")

	switch {
	case e.Person != "":
		printReportItem("  Person:", personName(a, e.Person))
		printReportItem("    id:", e.Person)
	case e.Place != "":
		printReportItem("  Place:", placeName(a, e.Place))
		printReportItem("    id:", e.Place)
	default:
		printReportItem("FOO:", "relationship or event")
	}
}

func placeName(a *glx.GLXFile, id string) string {
	if id == "" {
		return unspecifiedValue
	}

	p, ok := a.Places[id]
	if !ok {
		return unknown(id, "place")
	}

	if p.Name == "" {
		return unnamed(id, "place")
	}

	if p.ParentID == "" {
		return p.Name
	}

	return p.Name + ", " + placeName(a, p.ParentID)
}

func mediaTitle(a *glx.GLXFile, id string) string {
	if id == "" {
		return unspecifiedValue
	}

	m, ok := a.Media[id]
	if !ok {
		return unknown(id, "media")
	}

	if m.Title == "" {
		return unnamed(id, "media")
	}

	return m.Title
}

func personName(a *glx.GLXFile, id string) string {
	p, ok := a.Persons[id]
	if !ok {
		return unknown(id, "person")
	}

	name := glx.PersonDisplayName(p)
	if name == "" {
		return unnamed(id, "person")
	}

	return name
}

func repositoryName(a *glx.GLXFile, id string) string {
	if id == "" {
		return unspecifiedValue
	}

	p, ok := a.Repositories[id]
	if !ok {
		return unknown(id, "repository")
	}

	if p.Name == "" {
		return unnamed(id, "repository")
	}

	return p.Name
}

func sourceTitle(a *glx.GLXFile, id string) string {
	if id == "" {
		return unspecifiedValue
	}

	p, ok := a.Sources[id]
	if !ok {
		return unknown(id, "source")
	}

	if p.Title == "" {
		return unnamed(id, "source")
	}

	return p.Title
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

func printReference(label, id, value string) {
	printReportItem(label, value)
	printReportItem("  id:", id)
}

func printReportHeader(typ, title string) {
	fmt.Printf("=== %s: %s ===\n", typ, title)
}

func printReportItem(label string, value string) {
	if value == "" {
		value = unspecifiedValue
	}
	printReportLine(label, value)
}

func printReportLine(label, value string) {
	fmt.Printf("  %-18s%s\n", label, value)
}

func printSectionHeader(title string) {
	const width = 50
	prefix := "── " + title + " "
	remaining := max(width-utf8.RuneCountInString(prefix), 2)

	fmt.Println(prefix + strings.Repeat("─", remaining))
}

const unspecifiedValue = "—"

func unknown(id, typ string) string {
	return formattedLabeledID("unknown", id, typ)
}

func unnamed(id, typ string) string {
	return formattedLabeledID("unnamed", id, typ)
}

func formattedLabeledID(label, id, typ string) string {
	return fmt.Sprintf("%s %s id %s", label, typ, id)
}
