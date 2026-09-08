package describe

import (
	"fmt"

	"github.com/genealogix/glx/go-glx"
)

type Relationship Entity[glx.Relationship]

func (r *Relationship) Describe(rpt Report) {
	rpt.Begin("Relationship", r.id)

	rpt.Item("Type", r.entity.Type)

	for _, p := range r.Participants() {
		p.DescribeAsReference(rpt)
	}

	r.DescribeEvents(rpt)

	rpt.End()
}

func (r *Relationship) DescribeEvents(rpt Report) {
	r.EndEvent().DescribeAsSection(rpt, "Start Event")
	r.EndEvent().DescribeAsSection(rpt, "End Event")
}

func (r *Relationship) DescribeAsSection(rpt Report, label string) {
	rpt.BeginSection(fmt.Sprintf("%s: %s", label, r.id))

	rpt.Item("Type", r.entity.Type)

	for _, p := range r.Participants() {
		p.DescribeAsReference(rpt)
	}
}

func (r *Relationship) EndEvent() *Event {
	return nil
}

func (r *Relationship) StartEvent() *Event {
	return nil
}

func (r *Relationship) Participants() []*Participant {
	return r.archive.Participants(r.entity.Participants)
}
