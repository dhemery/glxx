package describe

import (
	"github.com/genealogix/glx/go-glx"
)

type Place Entity[glx.Place]

func (p *Place) Describe(r Report) {
	r.Begin("Place", p.id)

	r.Item("Name", p.Name())

	r.End()
}

func (p *Place) DescribeAsReference(r Report) {
	if p == nil {
		r.Item("Place", "")
		return
	}
	r.Reference("Place", p.Name(), p.id)
}

func (p *Place) DescribeAsSection(r Report, label string) {
}

func (p *Place) Name() string {
	name := p.entity.Name

	if name == "" {
		return name
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
