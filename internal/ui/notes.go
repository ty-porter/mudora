package ui

import (
	"fmt"

	. "modernc.org/tk9.0"
)

const (
	notesWidth  = 24 // characters
	notesHeight = 10 // lines
)

// Notes is a free-form scratchpad panel.
type Notes struct {
	text *TextWidget
}

func NewNotes() *Notes {
	opts := []Opt{Width(notesWidth), Height(notesHeight), Wrap("word"),
		Borderwidth(0), Highlightthickness(0), Padx(4), Pady(4)}
	// The classic text widget doesn't follow the ttk theme; color it by
	// hand, slightly lighter than the frame so the writable area reads.
	if fg := StyleLookup("TLabel", Foreground); fg != "" {
		opts = append(opts, Foreground(fg), Insertbackground(fg))
	}
	if bg := StyleLookup("TFrame", Background); bg != "" {
		opts = append(opts, Background(lighten(bg, 16)))
	}
	return &Notes{text: Text(opts...)}
}

// Widget returns the notes widget for placement in a parent layout.
func (n *Notes) Widget() *TextWidget { return n.text }

// lighten raises each channel of a #rrggbb color by the given amount,
// returning the input unchanged if it isn't in that form.
func lighten(hex string, by int) string {
	var r, g, b int
	if _, err := fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b); err != nil {
		return hex
	}
	return fmt.Sprintf("#%02x%02x%02x", min(r+by, 255), min(g+by, 255), min(b+by, 255))
}
