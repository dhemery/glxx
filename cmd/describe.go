package cmd

import (
	"fmt"

	"github.com/dhemery/glxx/load"
	"github.com/genealogix/glx/go-glx"
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe an entity",
	Long:  "Describe an entity",
	RunE:  describe,
	Args:  cobra.ExactArgs(1),
}

func describe(_ *cobra.Command, ids []string) error {
	archive, err := load.Load(archivePath)
	if err != nil {
		return err
	}

	id := ids[0]
	if _, ok := archive.Assertions[id]; ok {
		return fmt.Errorf("Not implemented: describe assertions")
	}
	if _, ok := archive.Citations[id]; ok {
		return fmt.Errorf("Not implemented: describe citations")
	}
	if e, ok := archive.Events[id]; ok {
		return describeEvent(archive, e)
	}
	if _, ok := archive.Media[id]; ok {
		return fmt.Errorf("Not implemented: describe media")
	}
	if _, ok := archive.Persons[id]; ok {
		return fmt.Errorf("Not implemented: describe person")
	}
	if _, ok := archive.Places[id]; ok {
		return fmt.Errorf("Not implemented: describe place")
	}
	if _, ok := archive.Relationships[id]; ok {
		return fmt.Errorf("Not implemented: describe relationship")
	}
	if _, ok := archive.Repositories[id]; ok {
		return fmt.Errorf("Not implemented: describe repository")
	}
	if _, ok := archive.Sources[id]; ok {
		return fmt.Errorf("Not implemented: describe source")
	}

	return fmt.Errorf("Unknown ID: %s", id)
}

func describeEvent(a *glx.GLXFile, e *glx.Event) error {
	fmt.Println("      Title : ", e.Title)
	fmt.Println("       Type : ", e.Type)
	fmt.Println("      Place : ", fullPlaceName(a, e.PlaceID))
	fmt.Println("       Date : ", e.Date)
	fmt.Println("Participants:")
	for _, p := range e.Participants {
		fmt.Printf("    %s: %s\n", p.Role, p.Person)
	}
	return nil
}

func fullPlaceName(a *glx.GLXFile, id string) string {
	p, ok := a.Places[id]
	if !ok {
		return id
	}

	if p.Name == "" {
		return "nameless:" + id
	}

	if p.ParentID == "" {
		return p.Name
	}

	return p.Name + ", " + fullPlaceName(a, p.ParentID)
}
