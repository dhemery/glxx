package cmd

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"os"

	"github.com/dhemery/glxx/load"
	"github.com/genealogix/glx/go-glx"
	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump a GENEALOGIX archive as JSON.",
	Long:  "Dump a GENEALOGIX archive as JSON.",
	Args:  cobra.ExactArgs(0),
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

func dump(c *cobra.Command, args []string) error {
	archive, err := load.Load(archivePath)
	if err != nil {
		return err
	}

	var appliedFilters bool
	filteredArchive := new(glx.GLXFile)

	if dumpAssertions {
		filteredArchive.Assertions = archive.Assertions
		appliedFilters = true
	}

	if dumpCitations {
		filteredArchive.Citations = archive.Citations
		appliedFilters = true
	}

	if dumpEvents {
		filteredArchive.Events = archive.Events
		appliedFilters = true
	}

	if dumpMedia {
		filteredArchive.Media = archive.Media
		appliedFilters = true
	}

	if dumpPersons {
		filteredArchive.Persons = archive.Persons
		appliedFilters = true
	}

	if dumpPlaces {
		filteredArchive.Places = archive.Places
		appliedFilters = true
	}

	if dumpRelationships {
		filteredArchive.Relationships = archive.Relationships
		appliedFilters = true
	}

	if dumpRepositories {
		filteredArchive.Repositories = archive.Repositories
		appliedFilters = true
	}

	if dumpSources {
		filteredArchive.Sources = archive.Sources
		appliedFilters = true
	}

	if appliedFilters {
		archive = filteredArchive
	}

	return json.MarshalWrite(os.Stdout, archive, jsontext.WithIndent("  "), json.OmitZeroStructFields(true))
}
