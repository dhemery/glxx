package cmd

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"os"
	"strings"

	"github.com/dhemery/glxx/load"
	"github.com/genealogix/glx/go-glx"
	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump a GENEALOGIX archive as JSON.",
	Long:  "Dump a GENEALOGIX archive as JSON.",
	RunE:  dump,
}

var (
	dumpAssertions    = false
	dumpCitations     = false
	dumpEvents        = false
	dumpMedia         = false
	dumpPersons       = false
	dumpPlaces        = false
	dumpRelationships = false
	dumpRepositories  = false
	dumpSources       = false
)

func init() {
	dumpCmd.Flags().BoolVar(&dumpAssertions, "assertions", dumpAssertions, "Dump assertions")
	dumpCmd.Flags().BoolVar(&dumpCitations, "citations", dumpCitations, "Dump citations")
	dumpCmd.Flags().BoolVar(&dumpEvents, "events", dumpEvents, "Dump events")
	dumpCmd.Flags().BoolVar(&dumpMedia, "media", dumpMedia, "Dump media")
	dumpCmd.Flags().BoolVar(&dumpPersons, "persons", dumpPersons, "Dump persons")
	dumpCmd.Flags().BoolVar(&dumpPlaces, "places", dumpPlaces, "Dump places")
	dumpCmd.Flags().BoolVar(&dumpRelationships, "relationships", dumpRelationships, "Dump relationships")
	dumpCmd.Flags().BoolVar(&dumpRepositories, "repositories", dumpRepositories, "Dump repositories")
	dumpCmd.Flags().BoolVar(&dumpSources, "sources", dumpSources, "Dump sources")
}

func dump(c *cobra.Command, entityIDs []string) error {
	archive, err := load.Load(archivePath)
	if err != nil {
		return err
	}

	var selectedEntities bool
	filteredArchive := new(glx.GLXFile)

	if dumpAssertions {
		filteredArchive.Assertions = archive.Assertions
		selectedEntities = true
	}

	if dumpCitations {
		filteredArchive.Citations = archive.Citations
		selectedEntities = true
	}

	if dumpEvents {
		filteredArchive.Events = archive.Events
		selectedEntities = true
	}

	if dumpMedia {
		filteredArchive.Media = archive.Media
		selectedEntities = true
	}

	if dumpPersons {
		filteredArchive.Persons = archive.Persons
		selectedEntities = true
	}

	if dumpPlaces {
		filteredArchive.Places = archive.Places
		selectedEntities = true
	}

	if dumpRelationships {
		filteredArchive.Relationships = archive.Relationships
		selectedEntities = true
	}

	if dumpRepositories {
		filteredArchive.Repositories = archive.Repositories
		selectedEntities = true
	}

	if dumpSources {
		filteredArchive.Sources = archive.Sources
		selectedEntities = true
	}

	if len(entityIDs) > 0 {
		unknowns := oollectEntities(archive, filteredArchive, entityIDs)
		if len(unknowns) > 0 {
			return fmt.Errorf("Unknown IDs: %s", strings.Join(unknowns, ", "))
		}
		selectedEntities = true
	}

	if selectedEntities {
		archive = filteredArchive
	}

	return json.MarshalWrite(os.Stdout, archive,
		jsontext.WithIndent("  "),
		json.OmitZeroStructFields(true))
}

func oollectEntities(in *glx.GLXFile, out *glx.GLXFile, ids []string) []string {
	var unknowns = []string{}

	for _, id := range ids {
		if assertions, ok := in.Assertions[id]; ok {
			if out.Assertions == nil {
				out.Assertions = map[string]*glx.Assertion{}
			}
			out.Assertions[id] = assertions
			continue
		}
		if citation, ok := in.Citations[id]; ok {
			if out.Citations == nil {
				out.Citations = map[string]*glx.Citation{}
			}
			out.Citations[id] = citation
			continue
		}
		if event, ok := in.Events[id]; ok {
			if out.Events == nil {
				out.Events = map[string]*glx.Event{}
			}
			out.Events[id] = event
			continue
		}
		if media, ok := in.Media[id]; ok {
			if out.Media == nil {
				out.Media = map[string]*glx.Media{}
			}
			out.Media[id] = media
			continue
		}
		if person, ok := in.Persons[id]; ok {
			if out.Persons == nil {
				out.Persons = map[string]*glx.Person{}
			}
			out.Persons[id] = person
			continue
		}
		if place, ok := in.Places[id]; ok {
			if out.Places == nil {
				out.Places = map[string]*glx.Place{}
			}
			out.Places[id] = place
			continue
		}
		if relationship, ok := in.Relationships[id]; ok {
			if out.Relationships == nil {
				out.Relationships = map[string]*glx.Relationship{}
			}
			out.Relationships[id] = relationship
			continue
		}
		if repository, ok := in.Repositories[id]; ok {
			if out.Repositories == nil {
				out.Repositories = map[string]*glx.Repository{}
			}
			out.Repositories[id] = repository
			continue
		}
		if source, ok := in.Sources[id]; ok {
			if out.Sources == nil {
				out.Sources = map[string]*glx.Source{}
			}
			out.Sources[id] = source
			continue
		}
		unknowns = append(unknowns, id)
	}
	return unknowns
}
