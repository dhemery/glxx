package describe

import (
	"github.com/genealogix/glx/go-glx"
)

type Media Entity[glx.Media]

func (m *Media) Describe(r Report) {
	r.Begin(m)

	r.Item("Title", m.entity.Title)
	r.Item("URI", m.entity.URI)
	r.Item("Type", m.entity.Type)
	r.Item("MimeType", m.entity.MimeType)
	r.Item("Hash", m.entity.Hash)
	r.Item("Date", m.entity.Date.String())
	r.Reference(m.Source())

	r.End()
}

func (m *Media) Label() string {
	return "Media"
}

func (m *Media) ID() string {
	return m.id
}

func (m *Media) Name() string {
	return m.Title()
}

func (m *Media) Title() string {
	title := m.entity.Title
	if title == "" {
		return unnamed(m.id, "media")
	}

	return title
}

func (m *Media) Source() *Source {
	e := m.entity
	if e.Source == "" {
		return nil
	}
	return m.archive.Source(e.Source)
}

func mediaTitle(a *glx.GLXFile, id string) string {
	if id == "" {
		return unspecifiedValue
	}

	m, ok := a.Media[id]
	if !ok {
		return unknown(id, "media")
	}

	if m.Title == "" {
		return unnamed(id, "media")
	}

	return m.Title
}
