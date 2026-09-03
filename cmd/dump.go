package cmd

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"os"

	"github.com/dhemery/glxx/load"
	"github.com/spf13/cobra"
)

func init() {
	glxxCmd.AddCommand(dumpCmd)
}

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump a GENEALOGIX archive as JSON.",
	Long:  "Dump a GENEALOGIX archive as JSON.",
	Args:  cobra.ExactArgs(0),
	RunE:  dump,
}

func dump(c *cobra.Command, args []string) error {
	archive, err := load.Load(archivePath)
	if err != nil {
		return err
	}

	return json.MarshalWrite(os.Stdout, archive, jsontext.WithIndent("  "), json.OmitZeroStructFields(true))
}
