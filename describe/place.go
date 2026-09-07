package describe

import (
	"fmt"
	"io"

	"github.com/genealogix/glx/go-glx"
)

func printPlaceReference(a *glx.GLXFile, label, id string) {
	p := Archive{a}.Place(id)
	printReference(label, id, p.Name())
}

type Place Entity[glx.Place]

func (p *Place) Describe(w io.Writer) {
	printReportHeader("Place", p.id)
	fmt.Fprintln(w)

	printReportItem("Name:", p.Name())

	fmt.Fprintln(w)
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
