package describe

import (
	"fmt"

	"github.com/genealogix/glx/go-glx"
)

func printPersonReference(label string, p *Person) {
	printReportItem(label, p.Name())
	printReportItem("  id:", p.id)
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
