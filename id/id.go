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

var Command = &cobra.Command{
	Use:   "id [flags] entity...",
	Short: "Show and recommend IDs for GLX entities",
	RunE:  id,
	Args:  cobra.MinimumNArgs(1),
}

var (
	showFixes   = false
	showMatches = false
	showUnables = false
)

func init() {
	Command.Flags().BoolVarP(&showFixes, "fixes", "f", showFixes, "Show glx command to fix each mismatching ID")
	Command.Flags().BoolVarP(&showMatches, "matches", "m", showMatches, "Show each already matching ID")
	Command.Flags().BoolVarP(&showUnables, "unable", "u", showUnables, "Show reason if unable to recommend")
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
				fmt.Fprintf(os.Stdout, "id %s: %s\n", id, err)
			}
			continue
		}

		if id == recommendedID {
			if showMatches {
				fmt.Fprintln(os.Stdout, "ok:", id)
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

	default:
		return "", fmt.Errorf("type %T not implemented", v)
	}
}

func recommendEventID(event *glx.Event) (string, error) {
	participants := event.Participants

	principalParticipants, relationshipParticipants := categorizeEventParticipants(participants)

	isPersonEvent := len(principalParticipants) > 0
	isRelationshipEvent := len(relationshipParticipants) > 0

	if !isPersonEvent && !isRelationshipEvent {
		return "", fmt.Errorf("event has neither principals nor relationship participants: %s", participants)
	}

	if isPersonEvent && isRelationshipEvent {
		return "", fmt.Errorf("event has both principals and relationship participants: %s", participants)
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
	recommendedID := glx.EntityID(glx.EntityIDPrefixEvent, strings.Join(parts, "-"))

	return recommendedID, nil
}

func recommendPersonID(person *glx.Person) (string, error) {
	nameProp := person.Properties["name"]
	given, surname := glx.ExtractNameFields(nameProp)

	if given == "" || given == "—" || surname == "" || surname == "—" {
		return "", errors.New("no recommendation: person has empty given name or surname")
	}

	given, _, _ = strings.Cut(given, " ")

	return glx.EntityID(glx.EntityIDPrefixPerson, given+"-"+surname), nil
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

var relationshipEventPrimaryRoles = []string{
	"bride",
	"groom",
	"spouse",
}

var personEventPrincipalRoles = []string{
	"principal",
	"subject",
}
