package describe

import (
	"github.com/genealogix/glx/go-glx"
)

type Assertion Entity[glx.Assertion]

func (a *Assertion) ID() string {
	return a.id
}

func (a *Assertion) Label() string {
	return "Assertion"
}

func (a *Assertion) Describe(r Report) {
	r.Item("Status:", a.entity.Status)

	r.BeginReferenceSection(a.Subject())

	r.BeginSection("Conclusion")
	r.Item("Property:", a.entity.Property)
	p := a.Participant()
	if p == nil {
		r.Item("Value", a.entity.Value)
	} else {
		r.Reference(p)
	}
	r.Item("Date:", a.entity.Date.String())
	r.Item("Confidence:", a.entity.Confidence)

	for _, c := range a.Citations() {
		r.BeginReferenceSection(c)
		c.Describe(r)
	}

	for _, m := range a.Media() {
		r.BeginReferenceSection(m)
		m.Describe(r)
	}
	for _, s := range a.Sources() {
		r.BeginReferenceSection(s)
		s.Describe(r)
	}
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
	p := participant(a.archive, *a.entity.Participant)
	return &p
}

func (a *Assertion) Citations() []*Citation {
	return a.archive.Citations(a.entity.Citations)
}

func (a *Assertion) Subject() NamedSubject {
	s := a.entity.Subject
	switch {
	case s.Event != "":
		return a.archive.Event(s.Event)
	case s.Person != "":
		return a.archive.Person(s.Person)
	case s.Place != "":
		return a.archive.Place(s.Place)
	case s.Relationship != "":
		return a.archive.Relationship(s.Relationship)
	}
	return nil
}
