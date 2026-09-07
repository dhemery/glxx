package describe

import (
	"fmt"
	"io"

	"github.com/genealogix/glx/go-glx"
)

func printPersonReference(label string, p *Person) {
	printReportItem(label, p.Name())
	printReportItem("  id:", p.id)
}

type Person Entity[glx.Person]

func (p *Person) Describe(w io.Writer) {
	printReportHeader("Person", p.id)
	fmt.Fprintln(w)

	printReportItem("Name:", p.Name())

	fmt.Fprintln(w)
}

func (p *Person) Name() string {
	name := glx.PersonDisplayName(p.entity)
	if name == "" {
		return unnamed(p.id, "person")
	}

	return name
}
