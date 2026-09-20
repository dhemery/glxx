package describe

import (
	"strings"
)

type Participant struct {
	person *Person
	role   string
}

func (p *Participant) DescribeAsReference(r Report) {
	role := p.role
	if role == "" {
		role = "subject"
	}
	label := strings.ToUpper(role[:1]) + role[1:]
	p.person.DescribeAsReference(r, label)
}

func (p Participant) Name() string {
	return p.person.Name()
}
