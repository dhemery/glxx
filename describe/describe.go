// Package describe implements the glxx describe command.
package describe

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/dhemery/glxx/load"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "describe",
	Short: "Describe an entity",
	Long:  "Describe an entity",
	RunE:  describe,
	Args:  cobra.ExactArgs(1),
}

type Entity[T any] struct {
	archive Archive
	id      string
	entity  *T
}

type Describer interface {
	Describe(io.Writer)
}

func describe(c *cobra.Command, ids []string) error {
	archivePath, err := c.Flags().GetString("archive")
	if err != nil {
		return err
	}

	glxfile, err := load.Load(archivePath)
	if err != nil {
		return err
	}

	archive := Archive{glxfile}
	id := ids[0]

	entity := archive.Find(id)
	if entity == nil {
		return fmt.Errorf("Unknown ID: %s", id)
	}

	entity.Describe(os.Stdout)
	return nil
}

func printReportHeader(typ, title string) {
	fmt.Printf("=== %s: %s ===\n", typ, title)
}

func printReportItem(label string, value string) {
	if value == "" {
		value = unspecifiedValue
	}
	printReportLine(label, value)
}

func printReportLine(label, value string) {
	fmt.Printf("  %-18s%s\n", label, value)
}

func printReference(label, id, value string) {
	printReportItem(label, value)
	if id == "" {
		return
	}
	printReportItem("  id:", id)
}

func printSectionHeader(title string) {
	const width = 50
	prefix := "── " + title + " "
	remaining := max(width-utf8.RuneCountInString(prefix), 2)

	fmt.Println(prefix + strings.Repeat("─", remaining))
}

const unspecifiedValue = "—"

func unknown(id, typ string) string {
	return formattedLabeledID("unknown", id, typ)
}

func unnamed(id, typ string) string {
	return formattedLabeledID("unnamed", id, typ)
}

func formattedLabeledID(label, id, typ string) string {
	return fmt.Sprintf("%s %s id %s", label, typ, id)
}
