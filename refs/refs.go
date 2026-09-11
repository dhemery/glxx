// Package refs implements the glxx refs command.
package refs

import (
	"fmt"
	"os"
	"slices"

	"github.com/dhemery/glxx/describe"
	"github.com/dhemery/glxx/load"
	"github.com/genealogix/glx/go-glx"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "refs [flags] entity",
	Short: "Describes each entity that references the given one",
	Args:  cobra.ExactArgs(1),
	RunE:  runRefs,
}

var (
	refsCount = false
	refsList  = false
)

func init() {
	Command.Flags().BoolVarP(&refsCount, "count", "c", refsCount, "Count referrers instead of describing")
	Command.Flags().BoolVarP(&refsList, "list", "l", refsList, "List referrer IDs instead of describing")
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

	id := args[0]

	referrerIDs := referrersTo(id, glxfile)

	if refsCount {
		fmt.Fprintln(os.Stdout, len(referrerIDs), "references to", id)
		return nil
	}

	slices.Sort(referrerIDs)

	if refsList {
		for _, referrerID := range referrerIDs {
			fmt.Fprintln(os.Stdout, referrerID)
		}
		return nil
	}

	r := describe.NewReport(os.Stdout)

	archive := describe.Archive{File: glxfile}

	for _, referrerID := range referrerIDs {
		referrer := archive.Find(referrerID)
		referrer.Describe(r)
	}

	return nil
}

func referrersTo(id string, f *glx.GLXFile) []string {
	var out []string

	for referrerID, referrer := range f.Assertions {
		if assertionRefersTo(referrer, id) {
			out = append(out, referrerID)
		}
	}

	return out

}

func assertionRefersTo(a *glx.Assertion, id string) bool {
	subject := a.Subject
	if subject.Event == id || subject.Person == id || subject.Place == id || subject.Relationship == id {
		return true
	}
	if a.Value == id {
		return true
	}
	if a.Participant != nil && a.Participant.Person == id {
		return true
	}
	if slices.Contains(a.Sources, id) {
		return true
	}
	if slices.Contains(a.Citations, id) {
		return true
	}
	if slices.Contains(a.Media, id) {
		return true
	}
	return false
}

func citationRefersTo(c *glx.Citation, id string) bool         { return false }
func eventRefersTo(e *glx.Event, id string) bool               { return false }
func mediaRefersTo(m *glx.Media, id string) bool               { return false }
func personRefersTo(p *glx.Person, id string) bool             { return false }
func placeRefersTo(p *glx.Place, id string) bool               { return false }
func relatiohshipRefersTo(r *glx.Relationship, id string) bool { return false }
func sourceRefersTo(s *glx.Source, id string) bool             { return false }
