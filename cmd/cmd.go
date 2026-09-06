// Package cmd defines the subcommands for glxx
package cmd

import (
	"fmt"
	"os"

	"github.com/dhemery/glxx/load"
	"github.com/genealogix/glx/go-glx"
	"github.com/spf13/cobra"
)

var glxxCmd = &cobra.Command{
	Use:               "glxx",
	Short:             "GLXX is Dale's companion to glx.",
	Long:              "GLXX is Dale's companion to glx.",
	Version:           "v0.0.0.unsupported.0",
	SilenceErrors:     true,
	PersistentPreRunE: loadArchive,
}

var archivePath string = "."
var archive *glx.GLXFile

func init() {
	if p := os.Getenv("GLXX_ARCHIVE"); p != "" {
		archivePath = p
	}
	glxxCmd.PersistentFlags().StringVarP(&archivePath, "archive", "a", archivePath, "the `dir` of the archive")

	glxxCmd.AddCommand(describeCmd)
	glxxCmd.AddCommand(dumpCmd)
}

// Execute runs the glxx command.
func Execute() {
	if err := glxxCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func loadArchive(cmd *cobra.Command, _ []string) error {
	// SilenceUsage (after arg validation) so that arg-count errors still
	// show usage but runtime errors from loadArchive and RunE do not.
	cmd.SilenceUsage = true
	a, err := load.Load(archivePath)
	if err != nil {
		return fmt.Errorf("loading archive: %w", err)
	}
	archive = a
	return nil
}
