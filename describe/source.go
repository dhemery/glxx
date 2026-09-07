package describe

import (
	"fmt"
	"io"

	"github.com/genealogix/glx/go-glx"
)

type Source Entity[glx.Source]

func (s *Source) Describe(w io.Writer) {
	entity := s.entity

	printReportHeader("Source", s.id)
	fmt.Fprintln(w)

	printReportItem("Title:", s.Title())
	for _, author := range entity.Authors {
		printReportItem("Author:", author)
	}
	printReportItem("Date:", entity.Date.String())
	printReportItem("Language:", entity.Language)

	s.Repository().PrintReference()

	for _, m := range s.Media() {
		m.PrintReference()
	}

	fmt.Fprintln(w)
}

func (s *Source) Media() []*Media {
	var out []*Media
	for _, m := range s.entity.Media {

		out = append(out, s.archive.Media(m))
	}
	return out
}

func (s *Source) Repository() *Repository {
	return s.archive.Repository(s.id)
}

func (s *Source) Title() string {
	title := s.entity.Title
	if title == "" {
		return unnamed(s.id, "source")
	}

	return title
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

func (s *Source) PrintReference() {
	printReference("Source:", s.id, s.Title())
}

func printSourceReference(a *glx.GLXFile, label, id string) {
	printReference(label, id, sourceTitle(a, id))
}

func (s *Source) PrintSection() {
	const header = "Source: %s %s"

	if s == nil {
		printSectionHeader(fmt.Sprintf(header, "", "(unknown)"))
		return
	}
	printSectionHeader(fmt.Sprintf(header, s.id, ""))
	printReportItem("Source:", s.Title())
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
