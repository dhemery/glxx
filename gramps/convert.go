package gramps

import (
	"fmt"
	"strings"

	"github.com/genealogix/glx/go-glx"
)

func importGramps(grampsData *DB) (*glx.GLXFile, error) {
	g, err := newGLX()
	if err != nil {
		return nil, err
	}

	for handle, in := range grampsData.Citations {
		id, out, err := convertCitation(handle, in, grampsData)
		if err != nil {
			return nil, err
		}
		g.Citations[id] = out
	}

	for handle, in := range grampsData.Events {
		id, out, err := convertEvent(handle, in, grampsData)
		if err != nil {
			return nil, err
		}
		g.Events[id] = out
	}

	for handle, in := range grampsData.Media {
		id, out, err := convertMedia(handle, in, grampsData)
		if err != nil {
			return nil, err
		}
		g.Media[id] = out
	}

	for handle, in := range grampsData.People {
		id, out, err := convertPerson(handle, in, grampsData)
		if err != nil {
			return nil, err
		}
		g.Persons[id] = out
	}

	for handle, in := range grampsData.Places {
		id, out, err := convertPlace(handle, in, grampsData)
		if err != nil {
			return nil, err
		}
		g.Places[id] = out
	}

	for handle, in := range grampsData.Repositories {
		id, out, err := convertRepository(handle, in, grampsData)
		if err != nil {
			return nil, err
		}
		g.Repositories[id] = out
	}

	for handle, in := range grampsData.Sources {
		id, out, err := convertSource(handle, in, grampsData)
		if err != nil {
			return nil, err
		}
		g.Sources[id] = out
	}

	return g, nil
}

func newGLX() (*glx.GLXFile, error) {
	g := new(glx.GLXFile)

	// if err := glx.LoadStandardVocabulariesIntoGLX(g); err != nil {
	// 	return nil, fmt.Errorf("loading standard vocabularies: %w", err)
	// }

	g.Assertions = make(map[string]*glx.Assertion)
	g.Citations = make(map[string]*glx.Citation)
	g.Events = make(map[string]*glx.Event)
	g.Media = make(map[string]*glx.Media)
	g.Persons = make(map[string]*glx.Person)
	g.Places = make(map[string]*glx.Place)
	g.Relationships = make(map[string]*glx.Relationship)
	g.Repositories = make(map[string]*glx.Repository)
	g.Sources = make(map[string]*glx.Source)

	return g, nil
}

func convertCitation(handle string, in *Citation, grampsData *DB) (string, *glx.Citation, error) {
	out := new(glx.Citation)

	out.Media = mediaRefIDs(in.Media, grampsData)

	return in.ID, out, nil
}

func convertEvent(handle string, in *Event, grampsData *DB) (string, *glx.Event, error) {
	out := new(glx.Event)

	date, err := convertDate(in.Date)
	if err != nil {
		return "", nil, fmt.Errorf("converting event %s %s date: %w", in.ID, handle, err)
	}
	out.Date = date

	return in.ID, out, nil
}

func convertMedia(handle string, in *Media, grampsData *DB) (string, *glx.Media, error) {
	out := new(glx.Media)

	file := in.File
	out.URI = file.Source
	out.Hash = file.Checksum
	out.Title = file.Description
	out.MimeType = file.Mime

	date, err := convertDate(in.Date)
	if err != nil {
		return "", nil, fmt.Errorf("converting media %s %s date: %w", in.ID, handle, err)
	}
	out.Date = glx.DateString(date)

	return in.ID, out, nil
}

func convertDate(in Date) (glx.DateString, error) {
	if in.DateStr.Val != "" {
		return convertDateStr(in.DateStr), nil
	}
	if in.DateVal.Val != "" {
		return convertDateVal(in.DateVal)
	}
	return convertDateSpan(in.DateSpan), nil
}

func convertDateSpan(in DateSpan) glx.DateString {

	var parts []string
	if in.Start != "" {
		parts = append(parts, "FROM", in.Start)
	}
	if in.Stop != "" {
		parts = append(parts, "TO", in.Start)
	}
	out := strings.Join(parts, " ")
	return glx.DateString(out)
}

func convertDateStr(in DateStr) glx.DateString {
	return glx.DateString(in.Val)
}

var dateValTypes = map[string]string{
	"about":  "ABT",
	"after":  "AFT",
	"before": "BEF",
}

func convertDateVal(in DateVal) (glx.DateString, error) {
	if in.Type == "" {
		return glx.DateString(in.Val), nil
	}

	qualifier, ok := dateValTypes[in.Type]
	if !ok {
		return "", fmt.Errorf("unknown DateVal.Type %s", in.Type)
	}

	date := fmt.Sprintf("%s %s", qualifier, in.Val)
	return glx.DateString(date), nil
}

func mediaRefIDs(rr []MediaRef, grampsData *DB) []string {
	var out []string
	for _, ref := range rr {
		id := mediaRefID(ref, grampsData)
		out = append(out, id)
	}
	return out
}

func mediaRefID(in MediaRef, grampsData *DB) string {
	h := in.MediaHandle
	m := grampsData.Media[h]
	if m == nil {
		return ""
	}
	return m.ID
}

func convertPerson(handle string, in *Person, grampsData *DB) (string, *glx.Person, error) {
	out := new(glx.Person)

	return in.ID, out, nil
}

func convertPlace(handle string, in *Place, grampsData *DB) (string, *glx.Place, error) {
	out := new(glx.Place)

	return in.ID, out, nil
}

func convertRepository(handle string, in *Repository, grampsData *DB) (string, *glx.Repository, error) {
	out := new(glx.Repository)

	return in.ID, out, nil
}

func convertSource(handle string, in *Source, grampsData *DB) (string, *glx.Source, error) {
	out := new(glx.Source)

	out.Media = mediaRefIDs(in.Media, grampsData)

	return in.ID, out, nil
}
