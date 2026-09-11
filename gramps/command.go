// Package gramps unmarshals Gramps data structures from XML.
package gramps

import (
	"os"

	"github.com/dhemery/glxx/dump"
	"github.com/spf13/cobra"
)

var (
	importDump = false
)

func init() {
	Command.Flags().BoolVar(&importDump, "dump", importDump, "Dump loaded Gramps XML to JSON and exit")
}

var Command = &cobra.Command{
	Use:   "import [flags] file",
	Short: "Import a Gramps XML file",
	Args:  cobra.ExactArgs(1),
	RunE:  run,
}

func run(c *cobra.Command, args []string) error {
	raw, err := loadGrampsXML(args[0])
	if err != nil {
		return err
	}

	if importDump {
		return dump.WriteJSON(os.Stdout, raw)
	}

	db := newDB(raw)
	converted, err := importGramps(db)
	if err != nil {
		return err
	}

	return dump.WriteJSON(os.Stdout, converted)
}
