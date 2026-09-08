// Package describe implements the glxx describe command.
package describe

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhemery/glxx/load"
	"github.com/spf13/cobra"
)

var describeLongUsage = `Describe an entity.

glxx describe describes one or more GLX entities. You can specify the
entities either by ID or by file path. Specifying the ID works in
single-file and multi-file archives. Specifying the file path works
only in multi-file archives where each entity's ID is used as its
file name.
`

var Command = &cobra.Command{
	Use:   "describe [flags] entity...",
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

func describe(c *cobra.Command, args []string) error {
	archivePath, err := c.Flags().GetString("archive")
	if err != nil {
		return err
	}

	glxfile, err := load.Load(archivePath)
	if err != nil {
		return err
	}

	archive := Archive{glxfile}
	r := Report{os.Stdout}

	var unknowns []string

	for _, arg := range args {
		id := strings.TrimSuffix(filepath.Base(arg), filepath.Ext(arg))
		entity := archive.Find(id)
		if entity == nil {
			unknowns = append(unknowns, arg)
			continue
		}

		entity.Describe(r)
	}

	if len(unknowns) > 0 {
		return fmt.Errorf("Unknown entities: %s", strings.Join(unknowns, ", "))
	}

	return nil
}
