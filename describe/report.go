package describe

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

type Report struct {
	w                    io.Writer
	IncludeNotes         bool
	IncludeAllProperties bool
}

func NewReport(w io.Writer) Report {
	return Report{w: w}
}

func (r Report) Begin(label, id string) {
	fmt.Fprintf(r.w, "=== %s: %s ===\n", label, id)
	fmt.Fprintln(r.w)
}

func (r Report) BeginSection(title string) {
	const width = 50
	prefix := "── " + title + " "
	remaining := max(width-utf8.RuneCountInString(prefix), 2)

	fmt.Fprintln(r.w)
	fmt.Fprintln(r.w, prefix+strings.Repeat("─", remaining))
}

func (r Report) End() {
	fmt.Fprintln(r.w)
}

func (r Report) Line(line ...any) {
	fmt.Fprintln(r.w, line...)
}

func (r Report) Item(label, value string) {
	const unspecifiedValue = "—"

	if value == "" {
		value = unspecifiedValue
	}
	fmt.Fprintf(r.w, "  %-18s%s\n", label[:min(len(label), 16)]+":", value)
}

func (r Report) Notes(nn []string) {
	if !r.IncludeNotes || len(nn) == 0 {
		return
	}
	for i, n := range nn {
		r.BeginSection(fmt.Sprintf("Note %d", i+1))
		fmt.Fprintln(r.w, n)
	}
}

func (r Report) Reference(label, value, id string) {
	r.Item(label, value)
	if id != "" {
		r.Item("  id", id)
	}
}
