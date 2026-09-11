// Package gramps unmarshals Gramps data structures from XML.
package gramps

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/dhemery/glxx/dump"
	"github.com/genealogix/glx/go-glx"
	"github.com/spf13/cobra"
)

var (
	importDump = ""
)

func init() {
	Command.Flags().StringVar(&importDump, "dump", importDump, "Dump gramps, glx, or files to JSON and exit")
}

var Command = &cobra.Command{
	Use:   "import [flags] file",
	Short: "Import a Gramps XML file",
	Args:  cobra.ExactArgs(1),
	RunE:  run,
}

var dumpTypes = []string{"gramps", "glx", "files"}

func run(c *cobra.Command, args []string) error {
	if importDump != "" && !slices.Contains(dumpTypes, importDump) {
		return fmt.Errorf("invalid dump type %q: expecting one of %s", importDump, dumpTypes)
	}
	raw, err := loadGrampsXML(args[0])
	if err != nil {
		return err
	}

	if importDump == "gramps" {
		return dump.WriteJSON(os.Stdout, raw)
	}

	db := newDB(raw)

	glxFile, err := importGramps(db)
	if err != nil {
		return err
	}
	if importDump == "glx" {
		return dump.WriteJSON(os.Stdout, glxFile)
	}

	s := glx.NewSerializer(nil)

	files, err := s.SerializeMultiFileToMap(glxFile)
	if err != nil {
		return err
	}

	if importDump == "files" {
		for f, c := range files {
			fmt.Fprintln(os.Stdout, f)
			fmt.Fprintln(os.Stdout, string(c))
		}
	}

	for f, c := range files {
		if err := os.MkdirAll(filepath.Dir(f), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(f, c, 0644); err != nil {
			return err
		}
	}
	return nil
}
