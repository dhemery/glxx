// Package gramps unmarshals Gramps data structures from XML.
package gramps

import (
	"fmt"
	"os"

	"github.com/dhemery/glxx/dump"
	"github.com/genealogix/glx/go-glx"
	"github.com/spf13/cobra"
)

var (
	importDumpGLX = false
	importDumpXML = false
)

func init() {
	Command.Flags().BoolVar(&importDumpXML, "dumpgramps", importDumpXML, "Dump loaded Gramps data to JSON and exit")
	Command.Flags().BoolVar(&importDumpGLX, "dumpglx", importDumpGLX, "Dump converted GLX to JSON and exit")
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

	if importDumpXML {
		return dump.WriteJSON(os.Stdout, raw)
	}

	db := newDB(raw)

	glxFile, err := importGramps(db)
	if err != nil {
		return err
	}
	if importDumpGLX {
		return dump.WriteJSON(os.Stdout, glxFile)
	}

	s := glx.NewSerializer(nil)

	serialized, err := s.SerializeMultiFileToMap(glxFile)
	if err != nil {
		return err
	}

	for f, c := range serialized {
		fmt.Fprintln(os.Stdout, f, string(c))
	}

	return nil
}
