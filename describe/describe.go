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
	return &Assertion{
		archive: a,
		id:      id,
		entity:  a.file.Assertions[id],
	}
}

func (a Archive) Citation(id string) *Citation {
	return &Citation{
		archive: a,
		id:      id,
		entity:  a.file.Citations[id],
	}
}

func (a Archive) Event(id string) *Event {
	return &Event{
		archive: a,
		id:      id,
		entity:  a.file.Events[id],
	}
}

func (a Archive) Media(id string) *Media {
	return &Media{
		archive: a,
		id:      id,
		entity:  a.file.Media[id],
	}
}

func (a Archive) Person(id string) *Person {
	return &Person{
		archive: a,
		id:      id,
		entity:  a.file.Persons[id],
	}
}

func (a Archive) Place(id string) *Place {
	return &Place{
		archive: a,
		id:      id,
		entity:  a.file.Places[id],
	}
}

func (a Archive) Relationship(id string) *Relationship {
	return &Relationship{
		archive: a,
		id:      id,
		entity:  a.file.Relationships[id],
	}
}

func (a Archive) Repository(id string) *Repository {
	return &Repository{
		archive: a,
		id:      id,
		entity:  a.file.Repositories[id],
	}
}

func (a Archive) Source(id string) *Source {
	return &Source{
		archive: a,
		id:      id,
		entity:  a.file.Sources[id],
	}
}

type Assertion Entity[glx.Assertion]

func (a *Assertion) Describe() {
	id := a.id
	glxfile := a.archive.file
	e := a.entity

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
		printParticipation(glxfile, p.Person, p.Role)
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
		printParticipation(a, p.Person, p.Role)
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
	id := p.id
	a := p.archive.file

	printReportHeader("Person", id)
	fmt.Println()

	printReportItem("Name:", personName(a, id))

	fmt.Println()
}

type Place Entity[glx.Place]

func (p *Place) Describe() {
	id := p.id
	a := p.archive.file

	printReportHeader("Place", id)
	fmt.Println()

	printReportItem("Name:", placeName(a, id))

	fmt.Println()
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
		printParticipation(a, p.Person, p.Role)
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
		printPersonSubjectSection(a, e.Person)
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
		printParticipation(a, p.Person, p.Role)
	}

	printRelationshipEvent(a, "Start", r.StartEvent)
	printRelationshipEvent(a, "End", r.EndEvent)
}

func printPersonSubjectSection(a *glx.GLXFile, id string) {
	printSectionHeader("Subject Person: " + id)
	printReportItem("Name:", personName(a, id))
}

func printPlaceSubjectSection(a *glx.GLXFile, id string) {
	printSectionHeader("Subject Place: " + id)
	printReportItem("Name:", placeName(a, id))
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
		printParticipation(a, p.Person, p.Role)
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
