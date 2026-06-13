package ui

import (
	"image"

	. "modernc.org/tk9.0"
)

const (
	// Sprites vary in size; each is centered in a fixed cell so grids stay
	// aligned. Sized to fit 16px sheet sprites at the icons package's
	// display scale.
	cellWidth  = 32
	cellHeight = 32
)

// iconGrid lays out icons in a fixed grid inside its own frame. Label and
// photo handles plus the icon source are kept so cells can be redrawn in
// place; a per-cell fingerprint lets unchanged cells skip the redraw.
type iconGrid struct {
	frame  *TFrameWidget
	labels []*LabelWidget
	// photos holds each cell's current Tk photo so it can be deleted when
	// replaced; otherwise Tk would orphan a photo image on every redraw.
	photos []*Img
	icon   func(int) (image.Image, error)
	// fingerprint reports a comparable summary of a cell's render-relevant
	// state. Refresh skips cells whose fingerprint hasn't changed. May be
	// nil, in which case every Refresh redraws.
	fingerprint func(int) any
	last        []any
}

// newIconGrid builds a grid of count cells, columns wide, fetching each
// cell's image from icon. Pass cellW/cellH of 0 to size cells to their
// images (the icon source must then produce uniformly sized images).
// columnMajor fills each column top-to-bottom before the next (cell 0 at the
// top-left, then straight down); the default fills each row left-to-right.
// fingerprint, if non-nil, lets Refresh skip cells whose state is unchanged;
// it must return a comparable value.
func newIconGrid(count, columns, cellW, cellH int, columnMajor bool, icon func(int) (image.Image, error), fingerprint func(int) any) (*iconGrid, error) {
	// Rows needed to hold count cells in the given number of columns; only
	// the column-major layout needs it.
	rows := (count + columns - 1) / columns
	g := &iconGrid{
		frame:       TFrame(),
		labels:      make([]*LabelWidget, 0, count),
		photos:      make([]*Img, 0, count),
		icon:        icon,
		fingerprint: fingerprint,
		last:        make([]any, count),
	}

	// Classic labels don't follow the ttk theme; match them to the themed
	// frame they sit in.
	bg := StyleLookup("TFrame", Background)

	for i := 0; i < count; i++ {
		img, err := icon(i)
		if err != nil {
			return nil, err
		}
		photo := NewPhoto(Data(img))
		opts := []Opt{Image(photo), Padx(0), Pady(0), Borderwidth(0)}
		if cellW > 0 {
			opts = append(opts, Width(cellW))
		}
		if cellH > 0 {
			opts = append(opts, Height(cellH))
		}
		if bg != "" {
			opts = append(opts, Background(bg))
		}
		label := g.frame.Label(opts...)
		row, col := i/columns, i%columns
		if columnMajor {
			row, col = i%rows, i/rows
		}
		Grid(label, Row(row), Column(col))
		g.labels = append(g.labels, label)
		g.photos = append(g.photos, photo)
		if fingerprint != nil {
			g.last[i] = fingerprint(i)
		}
	}

	return g, nil
}

// Frame returns the container widget for placement in a parent layout.
func (g *iconGrid) Frame() *TFrameWidget {
	return g.frame
}

// Refresh redraws cell i from its icon source. If a fingerprint is set and
// the cell's state is unchanged since the last draw, it's a no-op. The cell's
// previous photo is deleted so Tk doesn't accumulate orphaned images.
func (g *iconGrid) Refresh(i int) error {
	if g.fingerprint != nil {
		fp := g.fingerprint(i)
		if fp == g.last[i] {
			return nil
		}
		g.last[i] = fp
	}
	img, err := g.icon(i)
	if err != nil {
		return err
	}
	photo := NewPhoto(Data(img))
	g.labels[i].Configure(Image(photo))
	g.photos[i].Delete()
	g.photos[i] = photo
	return nil
}

// RefreshAll redraws every cell from its icon source, stopping at the first
// error.
func (g *iconGrid) RefreshAll() error {
	for i := range g.labels {
		if err := g.Refresh(i); err != nil {
			return err
		}
	}
	return nil
}

