package describe

import (
	"fmt"

	"github.com/genealogix/glx/go-glx"
)

type Citation Entity[glx.Citation]

func (c *Citation) Describe() {
	printReportHeader("Citation", c.id)
	fmt.Println()

	c.Source().PrintReference()
	c.Repository().PrintReference()

	for _, m := range c.Media() {
		m.PrintReference()
	}

	fmt.Println()
}

func (c *Citation) Media() []*Media {
	var out []*Media
	for _, m := range c.entity.Media {
		out = append(out, c.archive.Media(m))
	}
	return out
}

func (c *Citation) Source() *Source {
	return c.archive.Source(c.entity.SourceID)
}

func (c *Citation) Repository() *Repository {
	return c.archive.Repository(c.entity.RepositoryID)
}

func (c *Citation) PrintSection() {
	const header = "Citation: %s %s"
	printSectionHeader(fmt.Sprintf(header, c.id, ""))
	c.Source().PrintReference()
	c.Repository().PrintReference()
	for _, m := range c.Media() {
		m.PrintReference()
	}
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
