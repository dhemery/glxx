package describe

import (
	"github.com/genealogix/glx/go-glx"
)

type Citation Entity[glx.Citation]

func (c *Citation) Describe(r Report) {
	r.Reference(c.Source())
	r.Reference(c.Repository())

	for _, m := range c.Media() {
		r.Reference(m)
	}
}

func (c *Citation) ID() string {
	return c.id
}

func (c *Citation) Media() []*Media {
	var out []*Media
	for _, m := range c.entity.Media {
		out = append(out, c.archive.Media(m))
	}
	return out
}

func (c *Citation) Name() string {
	panic("unimplemented")
}

func (c *Citation) Label() string {
	return "Citation"
}

func (c *Citation) Repository() *Repository {
	return c.archive.Repository(c.entity.RepositoryID)
}

func (c *Citation) Source() *Source {
	return c.archive.Source(c.entity.SourceID)
}
