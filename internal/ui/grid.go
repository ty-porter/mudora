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

// iconGrid lays out icons in a fixed grid inside its own frame. Label
// handles and the icon source are kept so cells can be redrawn in place.
type iconGrid struct {
	frame  *TFrameWidget
	labels []*LabelWidget
	icon   func(int) (image.Image, error)
}

// newIconGrid builds a grid of count cells, columns wide, fetching each
// cell's image from icon. Pass cellW/cellH of 0 to size cells to their
// images (the icon source must then produce uniformly sized images).
func newIconGrid(count, columns, cellW, cellH int, icon func(int) (image.Image, error)) (*iconGrid, error) {
	g := &iconGrid{
		frame:  TFrame(),
		labels: make([]*LabelWidget, 0, count),
		icon:   icon,
	}

	// Classic labels don't follow the ttk theme; match them to the themed
	// frame they sit in.
	bg := StyleLookup("TFrame", Background)

	for i := 0; i < count; i++ {
		img, err := icon(i)
		if err != nil {
			return nil, err
		}
		opts := []Opt{Image(NewPhoto(Data(img))), Padx(0), Pady(0), Borderwidth(0)}
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
		Grid(label, Row(i/columns), Column(i%columns))
		g.labels = append(g.labels, label)
	}

	return g, nil
}

// Frame returns the container widget for placement in a parent layout.
func (g *iconGrid) Frame() *TFrameWidget {
	return g.frame
}

// Refresh redraws cell i from its icon source.
func (g *iconGrid) Refresh(i int) error {
	img, err := g.icon(i)
	if err != nil {
		return err
	}
	g.labels[i].Configure(Image(NewPhoto(Data(img))))
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

