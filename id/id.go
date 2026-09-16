// Package id implements the glxx id command
package id

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/dhemery/glxx/load"
	"github.com/genealogix/glx/go-glx"
	"github.com/spf13/cobra"
)

var idUsage = `Recommend entity IDs based on each GLX entity's type and fields.

glxx id applies opinionated rules to recommend IDs for GLX entities.
There are three possible results:
- If the ID matches the recommendation, indicate that.
- If the ID does not match the recommendation, show a 'glx rename'
  command to fix the ID.
- If glxx id has no recommendation, explain why.

By default glxx id shows all results. If any of the filter flags
('--fixes', '--matches', '--unable') are specified, only the selected
results are displayed.
`

var Command = &cobra.Command{
	Use:   "id [flags] entity...",
	Short: "Show and recommend IDs for GLX entities",
	Long:  idUsage,
	RunE:  id,
	Args:  cobra.MinimumNArgs(1),
}

var (
	showFixes   = false
	showMatches = false
	showUnables = false
)

func init() {
	Command.Flags().BoolVarP(&showFixes, "fixes", "f", showFixes, "Show the glx command to fix each mismatching ID")
	Command.Flags().BoolVarP(&showMatches, "matches", "m", showMatches, "Show each matching ID")
	Command.Flags().BoolVarP(&showUnables, "unable", "u", showUnables, "Show the reason if unable to recommend")
}

func id(c *cobra.Command, args []string) error {
	if !(showFixes || showMatches || showUnables) {
		showFixes = true
		showMatches = true
		showUnables = true
	}
	archivePath, err := c.Flags().GetString("archive")
	if err != nil {
		return err
	}

	glxFile, err := load.Load(archivePath)
	if err != nil {
		return err
	}

	for _, arg := range args {
		id := strings.TrimSuffix(filepath.Base(arg), filepath.Ext(arg))

		entity := findEntity(id, glxFile)
		if entity == nil {
			if id != arg {
				return fmt.Errorf("unknown entity %s (%s)", id, arg)
			}
			return fmt.Errorf("unknown entity %s", id)
		}

		recommendedID, err := recommendID(entity)
		if err != nil {
			if showUnables {
				fmt.Fprintf(os.Stdout, "no recommendation for %s: %s\n", id, err)
			}
			continue
		}

		// Accept the current ID if it starts with the recommended ID.
		// Assume that the suffix was added to differentiate this
		// entity from others with similar identifying attributes. So
		// a current ID person-george-washington-1732 would match
		// the recommended ID person-george-washington. The trailing
		// "-1732" differentiates this George Washington from others.
		if strings.HasPrefix(id, recommendedID) {
			if showMatches {
				fmt.Fprintln(os.Stdout, "matches recommendation:", id)
			}
			continue
		}

		if findEntity(recommendedID, glxFile) != nil {
			if showUnables {
				fmt.Fprintf(os.Stdout, "no recommendation for %s: recommended ID %s already in use\n",
					id, recommendedID)
			}
			continue
		}

		if showFixes {
			fmt.Fprintln(os.Stdout, "glx rename", id, recommendedID)
		}
	}

	return nil
}

func findEntity(id string, f *glx.GLXFile) any {
	if e, ok := f.Assertions[id]; ok {
		return e
	}
	if e, ok := f.Citations[id]; ok {
		return e
	}
	if e, ok := f.Events[id]; ok {
		return e
	}
	if e, ok := f.Media[id]; ok {
		return e
	}
	if e, ok := f.Persons[id]; ok {
		return e
	}
	if e, ok := f.Places[id]; ok {
		return e
	}
	if e, ok := f.Relationships[id]; ok {
		return e
	}
	if e, ok := f.Repositories[id]; ok {
		return e
	}
	if e, ok := f.Sources[id]; ok {
		return e
	}
	return nil

}

func recommendID(entity any) (string, error) {
	switch v := entity.(type) {
	case *glx.Event:
		return recommendEventID(v)

	case *glx.Person:
		return recommendPersonID(v)

	case *glx.Relationship:
		return recommendRelationshipID(v)

	default:
		return "", fmt.Errorf("entity type %T", v)
	}
}

func recommendEventID(e *glx.Event) (string, error) {
	summary := summarizeParticipants(e.Participants)

	slug, err := summary.slug(1)
	if err != nil {
		return "", fmt.Errorf("event: %w", err)
	}

	return fmt.Sprintf("%s%s-%s", glx.EntityIDPrefixEvent, slug, e.Type), nil
}

func recommendPersonID(person *glx.Person) (string, error) {
	nameProp := person.Properties["name"]
	given, surname := glx.ExtractNameFields(nameProp)

	if given == "" || given == "—" || surname == "" || surname == "—" {
		return "", errors.New("person has empty given name or surname")
	}

	given, _, _ = strings.Cut(given, " ")

	return glx.EntityID(glx.EntityIDPrefixPerson, given+"-"+surname), nil
}

func recommendRelationshipID(r *glx.Relationship) (string, error) {
	summary := summarizeParticipants(r.Participants)

	slug, err := summary.slug(2)
	if err != nil {
		return "", fmt.Errorf("relationship: %w", err)
	}

	return glx.EntityIDPrefixRelationship + slug, nil
}

// Roles that go on the left of a two-person ID.
var leftPrimaryRoles = []string{
	"groom",
	"parent",
	"adoptive_parent",
	"godparent",
}

// Roles that go on the right of a two-person ID.
var rightPrimaryRoles = []string{
	"bride",
	"child",
	"adopted_child",
	"godchild",
}

// Roles that go in source order in a two-person ID.
var unorderedPrimaryRoles = []string{
	"principal",
	"subject",
	"spouse",
	"sibling",
	"associate",
}

type participantSummary struct {
	Left      []glx.Participant
	Right     []glx.Participant
	Unordered []glx.Participant
}

func summarizeParticipants(pp []glx.Participant) participantSummary {
	var summary participantSummary

	for _, p := range pp {
		if slices.Contains(unorderedPrimaryRoles, p.Role) {
			summary.Unordered = append(summary.Unordered, p)
			continue
		}
		if slices.Contains(leftPrimaryRoles, p.Role) {
			summary.Left = append(summary.Left, p)
			continue
		}
		if slices.Contains(rightPrimaryRoles, p.Role) {
			summary.Right = append(summary.Right, p)
			continue
		}
	}

	return summary
}

func personSlug(personID string) string {
	return strings.TrimPrefix(personID, glx.EntityIDPrefixPerson)
}

func (s participantSummary) slug(minParticipants int) (string, error) {
	participantCount := len(s.Left) + len(s.Right) + len(s.Unordered)
	if participantCount < minParticipants {
		return "", fmt.Errorf("%d partipants with defining roles", participantCount)
	}

	if len(s.Unordered) == 2 {
		left := personSlug(s.Unordered[0].Person)
		right := personSlug(s.Unordered[1].Person)
		return fmt.Sprintf("%s-%s", left, right), nil
	}

	if len(s.Left) == 1 && len(s.Right) == 1 {
		left := personSlug(s.Left[0].Person)
		right := personSlug(s.Right[0].Person)
		return fmt.Sprintf("%s-%s", left, right), nil
	}

	if len(s.Unordered) == 1 {
		return personSlug(s.Unordered[0].Person), nil
	}

	var roles []string
	for _, p := range slices.Concat(s.Left, s.Right, s.Unordered) {
		roles = append(roles, p.Role)
	}

	slices.Sort(roles)

	return "", fmt.Errorf("incompatible roles: %s", roles)
}
