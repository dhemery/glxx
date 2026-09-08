// Package describe implements the glxx describe command.
package describe

import (
	"fmt"
	"os"
	"strings"

	"github.com/dhemery/glxx/load"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "describe [flags] id...",
	Short: "Describe an entity",
	Long:  "Describe an entity",
	RunE:  describe,
	Args:  cobra.MinimumNArgs(1),
}

type Entity[T any] struct {
	archive Archive
	id      string
	entity  *T
}

type Describer interface {
	Describe(Report)
}

func describe(c *cobra.Command, ids []string) error {
	archivePath, err := c.Flags().GetString("archive")
	if err != nil {
		return err
	}

	glxfile, err := load.Load(archivePath)
	if err != nil {
		return err
	}

	var unknowns []string
	r := Report{os.Stdout}

	archive := Archive{glxfile}
	for _, id := range ids {
		entity := archive.Find(id)
		if entity == nil {
			unknowns = append(unknowns, id)
			continue
		}

		entity.Describe(r)
	}

	if len(unknowns) > 0 {
		return fmt.Errorf("Unknown: %s", strings.Join(unknowns, ", "))
	}

	return nil
}
