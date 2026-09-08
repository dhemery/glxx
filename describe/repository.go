package describe

import (
	"github.com/genealogix/glx/go-glx"
)

type Repository Entity[glx.Repository]

func (r *Repository) Describe(rpt Report) {
	rpt.Begin("Repository", r.id)

	rpt.Item("Name", r.entity.Name)
	rpt.Item("Type", r.entity.Type)
	rpt.Item("Address", r.entity.Address)
	rpt.Item("City", r.entity.City)
	rpt.Item("State", r.entity.State)
	rpt.Item("Postal Code", r.entity.PostalCode)
	rpt.Item("Country", r.entity.Country)
	rpt.Item("Website", r.entity.Website)

	rpt.End()
}

func (r *Repository) DescribeAsReference(rpt Report) {
	if r == nil {
		rpt.Item("Repository", "")
		return
	}
	rpt.Reference("Repository", r.Name(), r.id)
}

func (r *Repository) Name() string {
	return r.entity.Name
}
