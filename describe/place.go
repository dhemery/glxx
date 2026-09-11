package describe

import (
	"fmt"

	"github.com/genealogix/glx/go-glx"
)

type Place Entity[glx.Place]

func (p *Place) Describe(r Report) {
	r.Begin("Place", p.id)

	r.Item("Name", p.entity.Name)
	p.Parent().DescribeAsReference(r, "Parent")
	r.Item("Type", p.entity.Type)
	r.Item("Latitude", angle(p.entity.Latitude))
	r.Item("Longitude", angle(p.entity.Longitude))

	r.Notes(p.entity.Notes)

	r.End()
}

func (p *Place) DescribeAsReference(r Report, label string) {
	if p == nil {
		r.Item(label, "")
		return
	}
	r.Reference(label, p.Name(), p.id)
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

func angle(f *float64) string {
	if f == nil {
		return ""
	}
	return fmt.Sprint(*f)
}
