package describe

import (
	"github.com/genealogix/glx/go-glx"
)

type Citation Entity[glx.Citation]

func (c *Citation) Describe(r Report) {
	r.Begin("Citation", c.id)

	c.Source().DescribeAsReference(r)
	c.Repository().DescribeAsReference(r)

	for _, m := range c.Media() {
		m.DescribeAsReference(r)
	}

	r.End()
}

func (c *Citation) DescribeAsSection(r Report) {
	if c == nil {
		return
	}
	r.BeginSection("Citation: " + c.id)

	c.Source().DescribeAsReference(r)
	c.Repository().DescribeAsReference(r)

	for _, m := range c.Media() {
		m.DescribeAsReference(r)
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

func (c *Citation) Repository() *Repository {
	return c.archive.Repository(c.entity.RepositoryID)
}

func (c *Citation) Source() *Source {
	return c.archive.Source(c.entity.SourceID)
}
