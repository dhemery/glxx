package describe

import (
	"github.com/genealogix/glx/go-glx"
)

type Relationship Entity[glx.Relationship]

// ID implements [NamedSubject].
func (r *Relationship) ID() string {
	panic("unimplemented")
}

// Name implements [NamedSubject].
func (r *Relationship) Name() string {
	panic("unimplemented")
}

// Label implements [NamedSubject].
func (r *Relationship) Label() string {
	panic("unimplemented")
}

func (r *Relationship) Describe(rpt Report) {

	rpt.Item("Type:", r.entity.Type)

	for _, p := range r.Participants() {
		rpt.Reference(p)
	}
}

func (r *Relationship) StartEvent() any {
	return nil
}

func (r *Relationship) Participants() []*Participant {
	return nil
}
