package describe

import (
	"github.com/genealogix/glx/go-glx"
)

type Source Entity[glx.Source]

func (s *Source) Describe(r Report) {
	entity := s.entity

	r.Begin(s)

	r.Item("Title:", s.Title())
	for _, author := range entity.Authors {
		r.Item("Author:", author)
	}
	r.Item("Date:", entity.Date.String())
	r.Item("Language:", entity.Language)

	r.Reference(s.Repository())

	for _, m := range s.Media() {
		r.Reference(m)
	}

	r.End()
}

func (_ *Source) Label() string {
	return "Source:"
}

func (s *Source) ID() string {
	return s.id
}

func (s *Source) Name() string {
	return s.Title()
}

func (s *Source) Media() []*Media {
	var out []*Media
	for _, m := range s.entity.Media {

		out = append(out, s.archive.Media(m))
	}
	return out
}

func (s *Source) Repository() *Repository {
	return s.archive.Repository(s.id)
}

func (s *Source) Title() string {
	title := s.entity.Title
	if title == "" {
		return unnamed(s.id, "source")
	}

	return title
}

func sourceTitle(a *glx.GLXFile, id string) string {
	if id == "" {
		return unspecifiedValue
	}

	p, ok := a.Sources[id]
	if !ok {
		return unknown(id, "source")
	}

	if p.Title == "" {
		return unnamed(id, "source")
	}

	return p.Title
}
