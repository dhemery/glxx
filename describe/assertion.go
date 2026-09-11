package describe

import (
	"github.com/genealogix/glx/go-glx"
)

type Assertion Entity[glx.Assertion]

func (a *Assertion) Describe(r Report) {
	r.Begin("Assertion", a.id)

	r.Item("Status", a.entity.Status)

	a.DescribeSubject(r)

	a.DescribeConclusion(r)

	for _, c := range a.Citations() {
		c.DescribeAsSection(r)
	}

	for _, m := range a.Media() {
		m.DescribeAsSection(r)
	}
	for _, s := range a.Sources() {
		s.DescribeAsSection(r)
	}

	r.Notes(a.entity.Notes)

	r.End()
}

func (a *Assertion) DescribeConclusion(r Report) {
	r.BeginSection("Conclusion")

	r.Item("Property", a.entity.Property)

	if p := a.Participant(); p != nil {
		p.DescribeAsReference(r)
	} else {
		r.Item("Value", a.entity.Value)
	}

	r.Item("Date", a.entity.Date.String())
	r.Item("Confidence", a.entity.Confidence)
}

func (a *Assertion) DescribeSubject(r Report) {
	s := a.entity.Subject
	switch {
	case s.Event != "":
		a.archive.Event(s.Event).DescribeAsSection(r, "Subject Event")
	case s.Person != "":
		a.archive.Person(s.Person).DescribeAsSection(r, "Subject Person")
	case s.Place != "":
		a.archive.Place(s.Place).DescribeAsSection(r, "Subject Place")
	case s.Relationship != "":
		a.archive.Relationship(s.Relationship).DescribeAsSection(r, "Subject Relationship")
	}
}

func (a *Assertion) Citations() []*Citation {
	return a.archive.Citations(a.entity.Citations)
}

func (a *Assertion) Sources() []*Source {
	return a.archive.Sources(a.entity.Sources)
}

func (a *Assertion) Media() []*Media {
	return a.archive.Medias(a.entity.Media)
}

func (a *Assertion) Participant() *Participant {
	if a.entity.Participant == nil {
		return nil
	}
	return a.archive.Participant(*a.entity.Participant)
}
