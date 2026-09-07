package describe

import (
	"github.com/genealogix/glx/go-glx"
)

type Place Entity[glx.Place]

func (p *Place) Describe(r Report) {
	r.Item("Name:", p.Name())
}

func (p *Place) ID() string {
	if p == nil {
		return ""
	}
	return p.id
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

func (p *Place) Label() string {
	return "Place"
}
