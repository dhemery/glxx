package describe

import (
	"github.com/genealogix/glx/go-glx"
)

type Source Entity[glx.Source]

func (s *Source) Describe(r Report) {
	r.Begin("Source", s.id)

	r.Item("Title", s.Title())
	for _, author := range s.entity.Authors {
		r.Item("Author", author)
	}
	r.Item("Date", s.entity.Date.String())
	r.Item("Language", s.entity.Language)

	s.Repository().DescribeAsReference(r)

	for _, m := range s.Media() {
		m.DescribeAsReference(r)
	}

	r.Notes(s.entity.Notes)

	r.End()
}

func (s *Source) DescribeAsReference(r Report) {
	if s == nil {
		r.Item("Source", "")
		return
	}
	r.Reference("Source", s.Title(), s.id)

}

func (s *Source) DescribeAsSection(r Report) {
	if s == nil {
		return
	}

	r.BeginSection("Source: " + s.id)

	r.Item("Title", s.Title())
	for _, author := range s.entity.Authors {
		r.Item("Author", author)
	}
	r.Item("Date", s.entity.Date.String())
	r.Item("Language", s.entity.Language)

	s.Repository().DescribeAsReference(r)

	for _, m := range s.Media() {
		m.DescribeAsReference(r)
	}
}

func (s *Source) Media() []*Media {
	var out []*Media
	for _, m := range s.entity.Media {

		out = append(out, s.archive.Media(m))
	}
	return out
}

func (s *Source) Repository() *Repository {
	return s.archive.Repository(s.entity.RepositoryID)
}

func (s *Source) Title() string {
	return s.entity.Title
}
