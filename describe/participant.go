package describe

import "github.com/genealogix/glx/go-glx"

type Participant struct {
	person *Person
	role   string
}

func (p Participant) ID() string {
	return p.person.ID()
}

func (p Participant) Name() string {
	return p.person.Name()
}

func (p Participant) Label() string {
	return p.role
}

func participant(a Archive, p glx.Participant) Participant {
	person := a.Person(p.Person)
	return Participant{person: person, role: p.Role}
}

func participants(a Archive, pp []glx.Participant) []Participant {
	var out []Participant
	for _, p := range pp {
		out = append(out, participant(a, p))
	}
	return out
}
