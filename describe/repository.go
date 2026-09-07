package describe

import (
	"github.com/genealogix/glx/go-glx"
)

type Repository Entity[glx.Repository]

func (r *Repository) Describe(rpt Report) {
	rpt.Item("Name:", r.entity.Name)
	rpt.Item("Type:", r.entity.Type)
	rpt.Item("Address:", r.entity.Address)
	rpt.Item("City:", r.entity.City)
	rpt.Item("State:", r.entity.State)
	rpt.Item("Postal Code:", r.entity.PostalCode)
	rpt.Item("Country:", r.entity.Country)
	rpt.Item("Website:", r.entity.Website)

}

func (r *Repository) ID() string {
	if r == nil {
		return ""
	}
	return r.id
}

func (r *Repository) Name() string {
	if r == nil {
		return unspecifiedValue
	}

	name := r.entity.Name

	if name == "" {
		return unnamed(r.id, "repository")
	}

	return name
}

func (r *Repository) Label() string {
	return "Repository"
}
