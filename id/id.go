// Package id implements the glxx id command
package id

import (
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

type result struct {
	arg           string
	id            string
	recommendedID string
	err           error
}

func id(c *cobra.Command, args []string) error {
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
			return fmt.Errorf("id %s: %w", id, err)
		}

		if id == recommendedID {
			fmt.Fprintln(os.Stdout, "good ID:", id)
		} else {
			fmt.Fprintln(os.Stdout, " bad ID:", recommendedID)
			fmt.Fprintln(os.Stdout, "    fix:", "glx", "rename", id, recommendedID)
		}
		fmt.Fprintln(os.Stdout)
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

	default:
		return "", fmt.Errorf("type %T not implemented", v)
	}
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

	if isPersonEvent {
		if len(principalParticipants) != 1 {
			return "", fmt.Errorf("event has %d principals: %s",
				len(principalParticipants), principalParticipants)
		}
		participantID := strings.TrimPrefix(principalParticipants[0].Person, glx.EntityIDPrefixPerson)
		recommendedID := glx.EntityIDPrefixEvent + participantID + "-" + event.Type
		return recommendedID, nil
	}

	if len(relationshipParticipants) != 2 {
		return "", fmt.Errorf("event has %d relationship participants: %s",
			len(relationshipParticipants), relationshipParticipants)
	}
	participant1ID := strings.TrimPrefix(relationshipParticipants[0].Person, glx.EntityIDPrefixPerson)
	participant2ID := strings.TrimPrefix(relationshipParticipants[1].Person, glx.EntityIDPrefixPerson)
	recommendedID := glx.EntityIDPrefixEvent + participant1ID + "-" + participant2ID + "-" + event.Type
	return recommendedID, nil
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
