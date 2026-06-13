package ui

import (
	"github.com/ty-porter/mudora/internal/ui/icons"
)

const itemGridColumns = 6

// ItemGrid displays the tracker's item icons.
type ItemGrid struct {
	*iconGrid
}

func NewItemGrid() (*ItemGrid, error) {
	// An item's icon is fully determined by its display tier, so the tier is
	// the redraw fingerprint.
	fingerprint := func(i int) any {
		it, err := icons.ItemAt(i)
		if err != nil {
			return nil
		}
		return it.State
	}
	g, err := newIconGrid(icons.Count(), itemGridColumns, cellWidth, cellHeight, false /* row-major */, icons.Item, fingerprint)
	if err != nil {
		return nil, err
	}
	return &ItemGrid{iconGrid: g}, nil
}
