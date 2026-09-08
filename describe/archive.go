package describe

import "github.com/genealogix/glx/go-glx"

type Archive struct {
	file *glx.GLXFile
}

func (a Archive) Find(id string) Describer {
	if d := a.Assertion(id); d != nil {
		return d
	}
	if d := a.Citation(id); d != nil {
		return d
	}
	if d := a.Event(id); d != nil {
		return d
	}
	if d := a.Media(id); d != nil {
		return d
	}
	if d := a.Person(id); d != nil {
		return d
	}
	if d := a.Place(id); d != nil {
		return d
	}
	if d := a.Relationship(id); d != nil {
		return d
	}
	if d := a.Repository(id); d != nil {
		return d
	}
	if d := a.Source(id); d != nil {
		return d
	}
	return nil
}

func (a Archive) Assertion(id string) *Assertion {
	entity := a.file.Assertions[id]
	if entity != nil {
		return &Assertion{a, id, entity}
	}
	return nil
}

func (a Archive) Citation(id string) *Citation {
	if entity, ok := a.file.Citations[id]; ok {
		return &Citation{a, id, entity}
	}
	return nil
}

func (a Archive) Citations(ids []string) []*Citation {
	var out []*Citation
	for _, id := range ids {
		out = append(out, a.Citation(id))
	}
	return out
}

func (a Archive) Event(id string) *Event {
	if entity, ok := a.file.Events[id]; ok {
		return &Event{a, id, entity}
	}
	return nil
}

func (a Archive) Media(id string) *Media {
	if entity, ok := a.file.Media[id]; ok {
		return &Media{a, id, entity}
	}
	return nil
}

func (a Archive) Medias(ids []string) []*Media {
	var out []*Media
	for _, id := range ids {
		out = append(out, a.Media(id))
	}
	return out
}

func (a Archive) Participant(p glx.Participant) *Participant {
	person := a.Person(p.Person)
	return &Participant{person: person, role: p.Role}
}

func (a Archive) Participants(pp []glx.Participant) []*Participant {
	var out []*Participant
	for _, p := range pp {
		out = append(out, a.Participant(p))
	}
	return out
}

func (a Archive) Person(id string) *Person {
	if entity, ok := a.file.Persons[id]; ok {
		return &Person{a, id, entity}
	}
	return nil
}

func (a Archive) Place(id string) *Place {
	if entity, ok := a.file.Places[id]; ok {
		return &Place{a, id, entity}
	}
	return nil
}

func (a Archive) Relationship(id string) *Relationship {
	if entity, ok := a.file.Relationships[id]; ok {
		return &Relationship{a, id, entity}
	}
	return nil
}

func (a Archive) Repository(id string) *Repository {
	if entity, ok := a.file.Repositories[id]; ok {
		return &Repository{a, id, entity}
	}
	return nil
}

func (a Archive) Source(id string) *Source {
	if entity, ok := a.file.Sources[id]; ok {
		return &Source{a, id, entity}
	}
	return nil
}

func (a Archive) Sources(ids []string) []*Source {
	var out []*Source
	for _, id := range ids {
		out = append(out, a.Source(id))
	}
	return out
}
