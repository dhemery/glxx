package describe

import (
	"github.com/genealogix/glx/go-glx"
)

type Person Entity[glx.Person]

func (p *Person) Describe(r Report) {
	r.Item("Name", p.Name())
}

func (p *Person) ID() string {
	return p.id
}

func (p *Person) Name() string {
	name := glx.PersonDisplayName(p.entity)
	if name == "" {
		return unnamed(p.id, p.Label())
	}

	return name
}

func (p *Person) Label() string {
	return "Person"
}
