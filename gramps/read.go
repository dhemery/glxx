package gramps

import (
	"encoding/xml"
	"io"
	"os"
)

func loadGrampsXML(fname string) (*rawGramps, error) {
	grampsFile, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer grampsFile.Close()

	xmlBytes, err := io.ReadAll(grampsFile)
	if err != nil {
		return nil, err
	}

	raw := new(rawGramps)
	err = xml.Unmarshal(xmlBytes, raw)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

type rawGramps struct {
	XMLName xml.Name `xml:"database"`
	Unknown
	Header
	Tags         []Tag        `xml:"tags>tag"`
	Citations    []Citation   `xml:"citations>citation"`
	Events       []Event      `xml:"events>event"`
	Families     []Family     `xml:"families>family"`
	Media        []Media      `xml:"objects>object"`
	Notes        []Note       `xml:"notes>note"`
	People       People       `xml:"people"`
	Places       []Place      `xml:"places>placeobj"`
	Repositories []Repository `xml:"repositories>repository"`
	Sources      []Source     `xml:"sources>source"`
}

type Header struct {
	XMLName    xml.Name   `xml:"header"`
	Created    Created    `xml:"created"`
	Researcher Researcher `xml:"researcher"`
	MediaPath  string     `xml:"mediapath"`
}

type Created struct {
	Date    string `xml:"date,attr"`
	Version string `xml:"version,attr"`
}

type Researcher struct {
	XMLName xml.Name `xml:"researcher"`
}

type People struct {
	Unknown
	People []Person `xml:"person"`

	HomePerson string `xml:"home,attr"`
}

// Unknown collects and reports unknown XML fields and attributes.
type Unknown struct {
	UnknownFields []UnknownField `xml:",any" json:",omitempty"`
	UnknownAttrs  []string       `xml:",any,attr" json:",omitempty"`
}

type UnknownField struct {
	XMLName xml.Name
	Value   string `xml:",innerxml"`
}

func (u Unknown) hasUnknowns() bool {
	return len(u.UnknownFields) > 0 || len(u.UnknownAttrs) > 0
}
