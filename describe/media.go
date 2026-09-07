package describe

import (
	"fmt"

	"github.com/genealogix/glx/go-glx"
)

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

func (m *Media) PrintReference() {
	printReference("Media:", m.id, m.Title())
}

func printMediaReference(a *glx.GLXFile, label, id string) {
	printReference(label, id, mediaTitle(a, id))
}

func (m *Media) Title() string {
	title := m.entity.Title
	if title == "" {
		return unnamed(m.id, "media")
	}

	return title
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
