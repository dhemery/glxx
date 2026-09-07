package describe

import (
	"fmt"

	"github.com/genealogix/glx/go-glx"
)

type Repository Entity[glx.Repository]

func (repository *Repository) Describe() {
	id := repository.id
	r := repository.entity

	printReportHeader("Repository", id)
	fmt.Println()

	printReportItem("Name:", r.Name)
	printReportItem("Type:", r.Type)
	printReportItem("Address:", r.Address)
	printReportItem("City:", r.City)
	printReportItem("State:", r.State)
	printReportItem("Postal Code:", r.PostalCode)
	printReportItem("Country:", r.Country)
	printReportItem("Website:", r.Website)

	fmt.Println()
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

func repositoryName(a *glx.GLXFile, id string) string {
	if id == "" {
		return unspecifiedValue
	}

	p, ok := a.Repositories[id]
	if !ok {
		return unknown(id, "repository")
	}

	if p.Name == "" {
		return unnamed(id, "repository")
	}

	return p.Name
}
func (r *Repository) PrintReference() {
	if r == nil {
		return
	}
	printReference("Repository:", r.Name(), r.id)
}

func printRepositoryReference(a *glx.GLXFile, label, id string) {
	name := repositoryName(a, id)
	if name == unspecifiedValue {
		printReportItem(label, name)
		return
	}
	printReference(label, id, name)
}
