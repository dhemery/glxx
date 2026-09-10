// Package gramps unmarshals Gramps data structures from XML.
package gramps

import (
	"errors"
	"os"

	"github.com/dhemery/glxx/dump"
	"github.com/spf13/cobra"
)

var (
	grampsImport = false
	grampsDump   = true
	grampsCheck  = false
)

var ErrHasUnknown = errors.New("has unknown fields or attrs")

func init() {
	Command.Flags().BoolVar(&grampsCheck, "check", grampsCheck, "Check for unknown fields and attributes in Gramps XML")
	Command.Flags().BoolVar(&grampsDump, "dump", grampsDump, "Dump Gramps XML to JSON")
	Command.Flags().BoolVar(&grampsImport, "import", grampsImport, "Import Gramps XML into GLX family archive")
}

var Command = &cobra.Command{
	Use:   "gramps [flags] file",
	Short: "Import a Gramps XML file or write as JSON",
	Args:  cobra.ExactArgs(1),
	RunE:  runGramps,
}

func runGramps(c *cobra.Command, args []string) error {
	if !(grampsImport || grampsDump) {
		c.Usage()
		return nil
	}

	raw, err := Read(args[0])
	if err != nil {
		return err
	}

	db := newDB(raw)

	if grampsDump {
		dump.WriteJSON(os.Stdout, db)
	}

	if grampsImport {
		return errors.ErrUnsupported
	}

	return nil
}
