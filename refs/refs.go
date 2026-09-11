// Package refs implements the glxx refs command.
package refs

import (
	"fmt"
	"os"

	"github.com/dhemery/glxx/describe"
	"github.com/dhemery/glxx/load"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "refs [flags] entity",
	Short: "Describes each entity that references the given one",
	Args:  cobra.ExactArgs(1),
	RunE:  runRefs,
}

func runRefs(c *cobra.Command, args []string) error {
	archivePath, err := c.Flags().GetString("archive")
	if err != nil {
		return err
	}

	glxfile, err := load.Load(archivePath)
	if err != nil {
		return err
	}

	archive := describe.Archive{File: glxfile}

	id := args[0]
	entity := archive.Find(id)
	if entity == nil {
		return fmt.Errorf("unknown entity: %s", id)
	}

	refs := findRefsTo(entity)
	if len(refs) == 0 {
		fmt.Fprintf(os.Stdout, "No references to %s\n", id)
		return nil
	}

	r := describe.NewReport(os.Stdout)

	for _, ref := range refs {
		ref.Describe(r)
	}

	return nil
}

func findRefsTo(e any) []describe.Describer {
	return nil
}
