// Package describe implements the glxx describe command.
package describe

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/dhemery/glxx/load"
	"github.com/genealogix/glx/go-glx"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "describe",
	Short: "Describe an entity",
	Long:  "Describe an entity",
	RunE:  describe,
	Args:  cobra.ExactArgs(1),
}

func describe(c *cobra.Command, ids []string) error {
	glxfile, err := load.Load(c)
	if err != nil {
		return err
	}

	archive := Archive{glxfile}
	id := ids[0]

	entity := archive.Find(id)
	if entity == nil {
		return fmt.Errorf("Unknown ID: %s", id)
	}

	entity.Describe()
	return nil
}

type Describer interface {
	Describe()
}

type Entity[T any] struct {
	archive Archive
	id      string
	entity  *T
}

type Archive struct {
	file *glx.GLXFile
}

func (a Archive) Find(id string) Describer {
	if d := a.Assertion(id); d != nil {
		return d
	}
	if d := a.Citation(id); d != nil {
		return d
	}
	if d := a.Event(id); d != nil {
		return d
	}
	if d := a.Media(id); d != nil {
		return d
	}
	if d := a.Person(id); d != nil {
		return d
	}
	if d := a.Place(id); d != nil {
		return d
	}
	if d := a.Relationship(id); d != nil {
		return d
	}
	if d := a.Repository(id); d != nil {
		return d
	}
	if d := a.Source(id); d != nil {
		return d
	}
	return nil
}

func (a Archive) Assertion(id string) *Assertion {
	entity := a.file.Assertions[id]
	if entity != nil {
		return &Assertion{a, id, entity}
	}
	return nil
}

func (a Archive) Citation(id string) *Citation {
	if entity, ok := a.file.Citations[id]; ok {
		return &Citation{a, id, entity}
	}
	return nil
}

func (a Archive) Event(id string) *Event {
	if entity, ok := a.file.Events[id]; ok {
		return &Event{a, id, entity}
	}
	return nil
}

func (a Archive) Media(id string) *Media {
	if entity, ok := a.file.Media[id]; ok {
		return &Media{a, id, entity}
	}
	return nil
}

func (a Archive) Person(id string) *Person {
	if entity, ok := a.file.Persons[id]; ok {
		return &Person{a, id, entity}
	}
	return nil
}

func (a Archive) Place(id string) *Place {
	if entity, ok := a.file.Places[id]; ok {
		return &Place{a, id, entity}
	}
	return nil
}

func (a Archive) Relationship(id string) *Relationship {
	if entity, ok := a.file.Relationships[id]; ok {
		return &Relationship{a, id, entity}
	}
	return nil
}

func (a Archive) Repository(id string) *Repository {
	if entity, ok := a.file.Repositories[id]; ok {
		return &Repository{a, id, entity}
	}
	return nil
}

func (a Archive) Source(id string) *Source {
	if entity, ok := a.file.Sources[id]; ok {
		return &Source{a, id, entity}
	}
	return nil
}

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

type Citation Entity[glx.Citation]

func (citation *Citation) Describe() {
	id := citation.id
	a := citation.archive.file
	c := citation.entity
	printReportHeader("Citation", id)
	fmt.Println()

	printSourceReference(a, "Source:", c.SourceID)
	printRepositoryReference(a, "Repository:", c.RepositoryID)
	for _, m := range c.Media {
		printMediaReference(a, "Media:", m)
	}

	fmt.Println()
}

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

type Media Entity[glx.Media]

func (media *Media) Describe() {
	m := media.entity
	id := media.id
	a := media.archive.file

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

type Person Entity[glx.Person]

func (p *Person) Describe() {
	printReportHeader("Person", p.id)
	fmt.Println()

	printReportItem("Name:", p.Name())

	fmt.Println()
}

func (p *Person) Name() string {
	name := glx.PersonDisplayName(p.entity)
	if name == "" {
		return unnamed(p.id, "person")
	}

	return name
}

type Place Entity[glx.Place]

func (p *Place) Describe() {
	printReportHeader("Place", p.id)
	fmt.Println()

	printReportItem("Name:", p.Name())

	fmt.Println()
}

func (p *Place) Name() string {
	if p == nil {
		return unspecifiedValue
	}

	name := p.entity.Name

	if name == "" {
		return unnamed(p.id, "place")
	}

	parent := p.Parent()
	if parent == nil {
		return name
	}

	return name + ", " + parent.Name()
}

func (p *Place) Parent() *Place {
	return p.archive.Place(p.entity.ParentID)
}

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

type Repository Entity[glx.Repository]

func (repository *Repository) Describe() {
	id := repository.id
	r := repository.entity

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

type Source Entity[glx.Source]

func (source *Source) Describe() {
	id := source.id
	s := source.entity
	a := source.archive.file

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

func printCitationSection(a *glx.GLXFile, id string) {
	const header = "Citation: %s %s"
	if id == "" {
		printSectionHeader(fmt.Sprintf(header, "(unspecified)", ""))
		return
	}

	c, ok := a.Citations[id]
	if !ok {
		printSectionHeader(fmt.Sprintf(header, id, "(unknown)"))
		return
	}

	printSectionHeader(fmt.Sprintf(header, id, ""))
	printReportItem("Source:", sourceTitle(a, c.SourceID))
	printRepositoryReference(a, "Repository:", c.RepositoryID)
	for _, m := range c.Media {
		printMediaReference(a, "Media", m)
	}
}

func printMediaReference(a *glx.GLXFile, label, id string) {
	printReference(label, id, mediaTitle(a, id))
}

func printPlaceReference(a *glx.GLXFile, label, id string) {
	p := Archive{a}.Place(id)
	printReference(label, id, p.Name())
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

func printSourceSection(a *glx.GLXFile, id string) {
	const header = "Source: %s %s"
	if id == "" {
		printSectionHeader(fmt.Sprintf(header, "(unspecified)", ""))
		return
	}

	_, ok := a.Sources[id]
	if !ok {
		printSectionHeader(fmt.Sprintf(header, id, "(unknown)"))
		return
	}
	printSectionHeader(fmt.Sprintf(header, id, ""))
	printReportItem("Source:", sourceTitle(a, id))
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

func printParticipation(p *Person, role string) {
	label := strings.ToUpper(role[:1]) + role[1:] + ":"
	printPersonReference(label, p)
}

func printPersonReference(label string, p *Person) {
	printReportItem(label, p.Name())
	printReportItem("  id:", p.id)
}

func printReference(label, id, value string) {
	printReportItem(label, value)
	if id == "" {
		return
	}
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
