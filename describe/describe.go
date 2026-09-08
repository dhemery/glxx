// Package describe implements the glxx describe command.
package describe

import (
	"fmt"
	"os"

	"github.com/dhemery/glxx/load"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "describe",
	Short: "Describe an entity",
	Long:  "Describe an entity",
	RunE:  describe,
	Args:  cobra.ExactArgs(1),
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

	archive := Archive{glxfile}
	id := ids[0]

	entity := archive.Find(id)
	if entity == nil {
		return fmt.Errorf("Unknown ID: %s", id)
	}

	r := Report{os.Stdout}

	entity.Describe(r)

	return nil
}
