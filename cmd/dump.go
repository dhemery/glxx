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

func dump(_ *cobra.Command, entityIDs []string) error {
	archiveIn, err := load.Load(archivePath)
	if err != nil {
		return err
	}

	var selectedEntities bool
	archiveOut := new(glx.GLXFile)

	if dumpAssertions {
		archiveOut.Assertions = archiveIn.Assertions
		selectedEntities = true
	}

	if dumpCitations {
		archiveOut.Citations = archiveIn.Citations
		selectedEntities = true
	}

	if dumpEvents {
		archiveOut.Events = archiveIn.Events
		selectedEntities = true
	}

	if dumpMedia {
		archiveOut.Media = archiveIn.Media
		selectedEntities = true
	}

	if dumpPersons {
		archiveOut.Persons = archiveIn.Persons
		selectedEntities = true
	}

	if dumpPlaces {
		archiveOut.Places = archiveIn.Places
		selectedEntities = true
	}

	if dumpRelationships {
		archiveOut.Relationships = archiveIn.Relationships
		selectedEntities = true
	}

	if dumpRepositories {
		archiveOut.Repositories = archiveIn.Repositories
		selectedEntities = true
	}

	if dumpSources {
		archiveOut.Sources = archiveIn.Sources
		selectedEntities = true
	}

	if len(entityIDs) > 0 {
		unknowns := oollectEntities(archiveIn, archiveOut, entityIDs)
		if len(unknowns) > 0 {
			return fmt.Errorf("Unknown IDs: %s", strings.Join(unknowns, ", "))
		}
		selectedEntities = true
	}

	if !selectedEntities {
		archiveOut = archiveIn
	}

	return json.MarshalWrite(os.Stdout, archiveOut,
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
