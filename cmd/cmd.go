// Package cmd defines the subcommands for glxx
package cmd

import (
	"fmt"
	"os"

	"github.com/dhemery/glxx/describe"
	"github.com/dhemery/glxx/dump"
	"github.com/dhemery/glxx/refs"
	"github.com/spf13/cobra"
)

var glxxCmd = &cobra.Command{
	Use:           "glxx",
	Short:         "GLXX is Dale's companion to glx.",
	Long:          "GLXX is Dale's companion to glx.",
	Version:       "v0.0.0.unsupported.0",
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		// SilenceUsage (after arg validation) so that arg-count errors still
		// show usage but runtime errors from RunE do not.
		cmd.SilenceUsage = true
	},
}

func init() {
	var archivePath string = "."
	if p := os.Getenv("GLXX_ARCHIVE"); p != "" {
		archivePath = p
	}
	glxxCmd.PersistentFlags().StringVarP(&archivePath, "archive", "a", archivePath, "the `dir` of the archive")

	glxxCmd.AddCommand(describe.Command)
	glxxCmd.AddCommand(dump.Command)
	glxxCmd.AddCommand(refs.Command)
}

// Execute runs the glxx command.
func Execute() {
	if err := glxxCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
