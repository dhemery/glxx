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

var describeLongUsage = `Describe GLX entities.

glxx describe describes one or more GLX entities.

You can specify entities either by ID or by file path.

NOTE: If you use a file path to specify an entity, glxx does not read
the file directly. Instead it removes the directory and extension
elements from the path and treats what remains as the entity ID. In
multi-file archives this allows you to use shell features and
commands to specify entities.

For example, to describe all assertions:

	glxx describe assertions/*.glx
`

var Command = &cobra.Command{
	Use:   "describe [flags] entity...",
	Short: "Describe GLX entities",
	Long:  describeLongUsage,
	RunE:  describe,
	Args:  cobra.MinimumNArgs(1),
}

var (
	describeNotes = false
)

func init() {
	Command.Flags().BoolVarP(&describeNotes, "notes", "n", describeNotes, "Include notes in each entity's description")
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
	r := Report{w: os.Stdout, IncludeNotes: describeNotes}

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
