package describe

import (
	"strings"
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
