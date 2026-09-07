package describe

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

type Report struct {
	w io.Writer
}

func (r Report) Begin(s Subject) {
	fmt.Fprintf(r.w, "=== %s: %s ===\n", s.Label(), s.ID())
}

func (r Report) BeginSection(title string) {
	const width = 50
	prefix := "── " + title + " "
	remaining := max(width-utf8.RuneCountInString(prefix), 2)

	fmt.Fprintln(r.w)
	fmt.Fprintln(r.w, prefix+strings.Repeat("─", remaining))
}

func (r Report) BeginReferenceSection(s NamedSubject) {
	title := fmt.Sprintf("%s: %s", s.Label(), s.ID())
	r.BeginSection(title)
}

func (r Report) Line(a ...any) {
	fmt.Fprintln(r.w, a...)
}

func (r Report) Item(label, value string) {
	if value == "" {
		value = unspecifiedValue
	}
	fmt.Fprintf(r.w, "  %-18s%s\n", label+":", value)
}

func (r Report) Reference(s NamedSubject) {
	r.Item(s.Label(), s.Name())
	id := s.ID()
	if id == "" {
		return
	}
	r.Item("  id:", s.ID())

}

func (r Report) End() {
	fmt.Fprintln(r.w)
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
