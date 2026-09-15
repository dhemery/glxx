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

func recommendEventID(event *glx.Event) (string, error) {
	if !slices.Contains(knownEventTypes, event.Type) {
		return "", fmt.Errorf("event type %s", event.Type)
	}

	participants := event.Participants

	principalParticipants, relationshipParticipants := categorizeEventParticipants(participants)

	isPersonEvent := len(principalParticipants) > 0
	isRelationshipEvent := len(relationshipParticipants) > 0

	if isPersonEvent && isRelationshipEvent {
		return "", fmt.Errorf("event has participants with both principal and relationship roles: %s",
			participants)
	}

	if !isPersonEvent && !isRelationshipEvent {
		return "", fmt.Errorf("event has no participants with principal or relationship roles: %s",
			participants)
	}

	var parts []string

	if isPersonEvent {
		if len(principalParticipants) != 1 {
			return "", fmt.Errorf("event has %d principals: %s",
				len(principalParticipants), principalParticipants)
		}
		parts = append(parts, strings.TrimPrefix(principalParticipants[0].Person, glx.EntityIDPrefixPerson))
	}

	if isRelationshipEvent {
		if len(relationshipParticipants) != 2 {
			return "", fmt.Errorf("event has %d relationship participants: %s",
				len(relationshipParticipants), relationshipParticipants)
		}
		parts = append(parts, strings.TrimPrefix(relationshipParticipants[0].Person, glx.EntityIDPrefixPerson))
		parts = append(parts, strings.TrimPrefix(relationshipParticipants[1].Person, glx.EntityIDPrefixPerson))
	}

	parts = append(parts, event.Type)
	return composeID(glx.EntityIDPrefixEvent, parts...), nil
}

var knownEventTypes = []string{
	"birth",
	"death",
	"marriage",
}

var relationshipEventPrimaryRoles = []string{
	"bride",
	"groom",
	"spouse",
}

var personEventPrincipalRoles = []string{
	"principal",
	"subject",
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
	if len(r.Participants) != 2 {
		return "", fmt.Errorf("relationship has %d participants", len(r.Participants))
	}

	byRole := participantsByRole(r.Participants)

	if len(byRole["parent"]) == 1 && len(byRole["child"]) == 1 {
		parent := strings.TrimPrefix(byRole["parent"][0].Person, glx.EntityIDPrefixPerson)
		child := strings.TrimPrefix(byRole["child"][0].Person, glx.EntityIDPrefixPerson)
		return composeID(glx.EntityIDPrefixRelationship, parent, child), nil
	}

	if len(byRole["spouse"]) == 2 {
		spouses := byRole["spouse"]
		s1 := strings.TrimPrefix(spouses[0].Person, glx.EntityIDPrefixPerson)
		s2 := strings.TrimPrefix(spouses[1].Person, glx.EntityIDPrefixPerson)
		return composeID(glx.EntityIDPrefixRelationship, s1, s2), nil
	}

	return "", fmt.Errorf("relationship has participants with roles: %s", r.Participants)
}

func categorizeEventParticipants(pp []glx.Participant) ([]glx.Participant, []glx.Participant) {
	var primaryParticipants []glx.Participant
	var relationshipParticipants []glx.Participant
	for _, p := range pp {
		if slices.Contains(personEventPrincipalRoles, p.Role) {
			primaryParticipants = append(primaryParticipants, p)
			continue
		}
		if slices.Contains(relationshipEventPrimaryRoles, p.Role) {
			relationshipParticipants = append(relationshipParticipants, p)
		}
	}
	return primaryParticipants, relationshipParticipants
}

func participantsByRole(pp []glx.Participant) map[string][]glx.Participant {
	var out = make(map[string][]glx.Participant)

	for _, p := range pp {
		out[p.Role] = append(out[p.Role], p)
	}

	return out
}

func composeID(prefix string, parts ...string) string {
	return prefix + strings.Join(parts, "-")
}
