package describe

import (
	"github.com/genealogix/glx/go-glx"
)

type Media Entity[glx.Media]

func (m *Media) Describe(r Report) {
	r.Begin("Media", m.id)

	r.Item("Title", m.entity.Title)
	r.Item("URI", m.entity.URI)
	r.Item("Type", m.entity.Type)
	r.Item("MimeType", m.entity.MimeType)
	r.Item("Hash", m.entity.Hash)
	r.Item("Date", m.entity.Date.String())

	m.properties.Describe(r)

	m.Source().DescribeAsReference(r)

	r.Notes(m.entity.Notes)

	r.End()
}

func (m *Media) DescribeAsReference(r Report) {
	if m == nil {
		r.Item("Media", "")
		return
	}
	r.Reference("Media", m.Title(), m.id)
}

func (m *Media) DescribeAsSection(r Report) {
	if m == nil {
		return
	}
	r.BeginSection("Media: " + m.id)
	r.Item("Title", m.entity.Title)
	r.Item("URI", m.entity.URI)
	r.Item("Type", m.entity.Type)
	r.Item("MimeType", m.entity.MimeType)
	r.Item("Hash", m.entity.Hash)
	r.Item("Date", m.entity.Date.String())

	m.Source().DescribeAsReference(r)

}

func (m *Media) Title() string {
	return m.entity.Title
}

func (m *Media) Source() *Source {
	e := m.entity
	if e.Source == "" {
		return nil
	}
	return m.archive.Source(e.Source)
}
