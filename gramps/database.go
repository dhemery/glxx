package gramps

// A DB is a map of the "primary" Gramps objects, indexed by their handles.
type DB struct {
	HomePerson   *Person
	Citations    map[string]*Citation
	Events       map[string]*Event
	Families     map[string]*Family
	Media        map[string]*Media
	Notes        map[string]*Note
	People       map[string]*Person
	Places       map[string]*Place
	Repositories map[string]*Repository
	Sources      map[string]*Source
}

func newDB(raw *rawGramps) *DB {
	g := &DB{
		Citations:    map[string]*Citation{},
		Events:       map[string]*Event{},
		Families:     map[string]*Family{},
		Media:        map[string]*Media{},
		Notes:        map[string]*Note{},
		People:       map[string]*Person{},
		Places:       map[string]*Place{},
		Repositories: map[string]*Repository{},
		Sources:      map[string]*Source{},
	}

	for _, item := range raw.Citations {
		g.Citations[item.Handle] = &item
	}

	for _, item := range raw.Events {
		g.Events[item.Handle] = &item
	}

	for _, item := range raw.Families {
		g.Families[item.Handle] = &item
	}

	for _, item := range raw.Media {
		g.Media[item.Handle] = &item
	}

	for _, item := range raw.Notes {
		g.Notes[item.Handle] = &item
	}

	for _, item := range raw.People.People {
		g.People[item.Handle] = &item
	}
	g.HomePerson = g.People[raw.People.HomePerson]

	for _, item := range raw.Places {
		g.Places[item.Handle] = &item
	}

	for _, item := range raw.Repositories {
		g.Repositories[item.Handle] = &item
	}

	for _, item := range raw.Sources {
		g.Sources[item.Handle] = &item
	}

	return g
}

// Record is a mix-in for the fields common to all records in a Gramps database.
type Record struct {
	Unknown
	// Identifes the record in the table.
	Handle string `xml:"handle,attr"`
	// The date and time of the latest update,
	// represented as a number of seconds since the Unix epoch.
	Change uint64 `xml:"change,attr"`
}

// PrimaryRecord is a mix-in for the fields common to all Gramps primary record types.
type PrimaryRecord struct {
	Record
	Privacy
	Tags
	// The object's Gramps ID.
	ID string `xml:"id,attr"`
}

type Citation struct {
	PrimaryRecord
	Date
	Page       string      `xml:"page"`
	Confidence uint8       `xml:"confidence"`
	Attributes []Attribute `xml:"srcattribute"`
	Media      []MediaRef  `xml:"objref"`
	Notes      []NoteRef   `xml:"noteref"`
	Sources    SourceRef   `xml:"sourceref"`
}

type Event struct {
	PrimaryRecord
	Date
	Type        string        `xml:"type"`
	Place       PlaceRef      `xml:"place"`
	Description string        `xml:"description"`
	Attributes  []Attribute   `xml:"attribute"`
	Citations   []CitationRef `xml:"citationref"`
	Notes       []NoteRef     `xml:"noteref"`
	Media       []MediaRef    `xml:"objref"`
}

type Family struct {
	PrimaryRecord
	Rel        FamilyType    `xml:"rel"`
	Father     PersonRef     `xml:"father"`
	Mother     PersonRef     `xml:"mother"`
	Children   []ChildRef    `xml:"childref"`
	Attributes []Attribute   `xml:"attribute"`
	Citations  []CitationRef `xml:"citationref"`
	Events     []EventRef    `xml:"eventref"`
	Media      []MediaRef    `xml:"objref"`
	Notes      []NoteRef     `xml:"noteref"`
}

type Media struct {
	PrimaryRecord
	Date
	File       MediaFile     `xml:"file"`
	Attributes []Attribute   `xml:"attribute"`
	Citations  []CitationRef `xml:"citationref"`
	Notes      []NoteRef     `xml:"noteref"`
}

type Note struct {
	PrimaryRecord
	Type   string      `xml:"type,attr"`
	Text   string      `xml:"text"`
	Styles []TextStyle `xml:"style"`
}

type Person struct {
	PrimaryRecord
	Privacy
	Gender        string         `xml:"gender"`
	Names         []PersonName   `xml:"name"`
	Addresses     []Address      `xml:"address"`
	Associations  []PersonRef    `xml:"personref"`
	ChildOf       []FamilyRef    `xml:"childof"`
	ParentIn      []FamilyRef    `xml:"parentin"`
	Attributes    []Attribute    `xml:"attribute"`
	Citations     []CitationRef  `xml:"citationref"`
	Events        []EventRef     `xml:"eventref"`
	LdsOrdinances []LDSOrdinance `xml:"lds_ord"`
	Media         []MediaRef     `xml:"objref"`
	Notes         []NoteRef      `xml:"noteref"`
	URLs          []URL          `xml:"url"`
}

type Place struct {
	PrimaryRecord
	Type          string        `xml:"type,attr"`
	Name          PlaceName     `xml:"pname"`
	Coordinates   Coordinates   `xml:"coord"`
	EncompassedBy []PlaceRef    `xml:"placeref"`
	Citations     []CitationRef `xml:"citationref"`
	Media         []MediaRef    `xml:"objref"`
	Notes         []NoteRef     `xml:"noteref"`
	URLs          []URL         `xml:"url"`
}

type Tag struct {
	Record
	Name     string `xml:"name,attr"`
	Color    string `xml:"color,attr"`
	Priority uint   `xml:"priority,attr"`
}

type Repository struct {
	PrimaryRecord
	Name      string    `xml:"rname"`
	Type      string    `xml:"type"`
	Addresses []Address `xml:"address"`
	Notes     []NoteRef `xml:"noteref"`
	URLs      []URL     `xml:"url"`
}

type Source struct {
	PrimaryRecord
	Title        string          `xml:"stitle"`
	Author       string          `xml:"sauthor"`
	PubInfo      string          `xml:"spubinfo"`
	Abbreviation string          `xml:"sabbrev"`
	Attributes   []Attribute     `xml:"srcattribute"`
	Media        []MediaRef      `xml:"objref"`
	Notes        []NoteRef       `xml:"noteref"`
	Repositories []RepositoryRef `xml:"reporef"`
}

type Address struct {
	Unknown `json:"unknown"`
	Date
	Street  string `xml:"street"`
	City    string `xml:"city"`
	State   string `xml:"state"`
	Country string `xml:"country"`
	Postal  string `xml:"postal"`
}

type Attribute struct {
	Unknown
	Type      string        `xml:"type,attr"`
	Value     string        `xml:"value,attr"`
	Citations []CitationRef `xml:"citationref"`
	Notes     []NoteRef     `xml:"noteref"`
}

type ChildRef struct {
	PersonRef
	FatherRelation string `xml:"frel,attr"`
	MotherRelation string `xml:"mrel,attr"`
}

type CitationRef struct {
	Unknown
	CitationHandle string `xml:"hlink,attr"`
}

type Coordinates struct {
	Unknown
	Longitude string `xml:"long,attr"`
	Latitude  string `xml:"lat,attr"`
}

type Date struct {
	Unknown
	DateSpan DateSpan `xml:"datespan"`
	DateStr  DateStr  `xml:"datestr"`
	DateVal  DateVal  `xml:"dateval"`
}

type DateStr struct {
	Val string `xml:"val,attr"`
}

// A DateVal represents a date or range of dates,
// possibly with qualifiers such as "about" or "before."
type DateVal struct {
	Val  string `xml:"val,attr"`
	Type string `xml:"type,attr"`
}

type DateSpan struct {
	Start string `xml:"start,attr"`
	Stop  string `xml:"stop,attr"`
}

// An EventRef is a references from a person to an event in which the person played a role.
type EventRef struct {
	Unknown
	EventHandle string        `xml:"hlink,attr"`
	Role        string        `xml:"role,attr"`
	Attributes  []Attribute   `xml:"attribute"`
	Citations   []CitationRef `xml:"citationref"`
	Notes       []NoteRef     `xml:"noteref"`
}

type FamilyRef struct {
	Unknown
	FamilyHandle string `xml:"hlink,attr"`
}

type FamilyType struct {
	Unknown
	Type string `xml:"type,attr"`
}

// LDSOrdinance matches the `lds_ord` element. I don't use it, so just match and ignore.
type LDSOrdinance struct {
	Unknown
}

type MediaFile struct {
	Unknown
	Source      string `xml:"src,attr"`
	Mime        string `xml:"mime,attr"`
	Description string `xml:"description,attr"`
	Checksum    string `xml:"checksum,attr"`
}

type MediaRef struct {
	Unknown
	MediaHandle string `xml:"hlink,attr"`
}
type NoteRef struct {
	Unknown
	NoteHandle string `xml:"hlink,attr"`
}

type PersonName struct {
	Unknown
	Privacy
	Date
	Type       string        `xml:"type,attr"`
	Alt        string        `xml:"alt,attr"`
	Title      string        `xml:"title"`
	First      string        `xml:"first"`
	Surname    string        `xml:"surname"`
	Suffix     string        `xml:"suffix"`
	Call       string        `xml:"call"`
	Nick       string        `xml:"nick"`
	FamilyNick string        `xml:"familynick"`
	Citations  []CitationRef `xml:"citationref"`
	Notes      []NoteRef     `xml:"noteref"`
}

type PersonRef struct {
	Unknown
	PersonHandle string        `xml:"hlink,attr"`
	Citations    []CitationRef `xml:"citationref"`
}

type PlaceName struct {
	Unknown
	Value string `xml:"value,attr"`
}

// PlaceRef reprents that a place is or was encompassed by an encompassing place.
type PlaceRef struct {
	Unknown
	Date
	// Reference to the encompassing place.
	PlaceHandle string `xml:"hlink,attr"`
	// The date or dates of the relationship.
}

// Privacy indicates whether an object is private
// and therefore should not be published.
type Privacy struct {
	Private bool `xml:"priv,attr"`
}

type RepositoryRef struct {
	Unknown
	RepositoryHandle string    `xml:"hlink,attr"`
	CallNumber       string    `xml:"callno,attr"`
	Medium           string    `xml:"medium,attr"`
	Notes            []NoteRef `xml:"noteref"`
}

type SourceRef struct {
	Unknown
	SourceHandle string `xml:"hlink,attr"`
}

type TextRange struct {
	Start int `xml:"start,attr"`
	End   int `xml:"end,attr"`
}

type TextStyle struct {
	Unknown
	Name   string      `xml:"name,attr"`
	Ranges []TextRange `xml:"range"`
}

type TagRef struct {
	Unknown
	TagHandle string `xml:"hlink,attr"`
}

// Tags is a mix-in that collects the containing element's
// references to tags.
type Tags struct {
	Tags []TagRef `xml:"tagref"`
}

type URL struct {
	Unknown
	HREF        string `xml:"href,attr"`
	Type        string `xml:"type,attr"`
	Description string `xml:"description,attr"`
}
