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
	g, err := newIconGrid(icons.Count(), itemGridColumns, cellWidth, cellHeight, icons.Item)
	if err != nil {
		return nil, err
	}
	return &ItemGrid{iconGrid: g}, nil
}
