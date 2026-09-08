package describe

import (
	"strings"

	"github.com/genealogix/glx/go-glx"
)

type Participant struct {
	person *Person
	role   string
}

func (p *Participant) DescribeAsReference(r Report) {
	label := strings.ToUpper(p.role[:1]) + p.role[1:]
	p.person.DescribeAsReference(r, label)
}

func (p Participant) Name() string {
	return p.person.Name()
}

func participant(a Archive, p glx.Participant) *Participant {
	person := a.Person(p.Person)
	return &Participant{person: person, role: p.Role}
}

func participants(a Archive, pp []glx.Participant) []*Participant {
	var out []*Participant
	for _, p := range pp {
		out = append(out, participant(a, p))
	}
	return out
}
