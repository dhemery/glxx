package describe

import (
	"fmt"

	"github.com/genealogix/glx/go-glx"
)

type Person Entity[glx.Person]

func (p *Person) Describe(r Report) {
	r.Begin("Person", p.id)

	r.Item("Name", p.Name())

	r.Notes(p.entity.Notes)

	r.End()
}

func (p *Person) DescribeAsReference(r Report, label string) {
	r.Reference(label, p.Name(), p.id)
}

func (p *Person) DescribeAsSection(r Report, label string) {
	if p == nil {
		return
	}
	r.BeginSection(fmt.Sprintf("%s: %s", label, p.id))
	r.Item("Name", p.Name())
}

func (p *Person) Name() string {
	return glx.PersonDisplayName(p.entity)
}
