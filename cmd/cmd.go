// Package cmd defines the subcommands for glxx
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var glxxCmd = &cobra.Command{
	Use:           "glxx",
	Short:         "GLXX is Dale's companion to glx.",
	Long:          "GLXX is Dale's companion to glx.",
	Version:       "v0.0.0.unsupported.0",
	SilenceErrors: true,
	// SilenceUsage is set in PersistentPreRun (after arg validation) so that
	// arg-count errors still show usage but runtime errors from RunE do not.
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		cmd.SilenceUsage = true
	},
}

var archivePath string = "."

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
