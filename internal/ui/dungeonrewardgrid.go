package ui

import (
	"image"

	"github.com/ty-porter/mudora/internal/ui/icons"
	. "modernc.org/tk9.0"
)

const rewardGridColumns = 2

// rewardLabels are the dungeon abbreviations, in dungeonRewards display
// order. They render in the pixel font to the left of each reward icon.
var rewardLabels = []string{"EP", "DP", "ToH", "PoD", "SP", "SW", "TT", "IP", "MM", "TR"}

// DungeonRewardGrid displays one cell per reward-bearing dungeon, showing
// the crystal or pendant it holds. Left click toggles whether the reward is
// acquired; right click cycles which reward the dungeon holds.
type DungeonRewardGrid struct {
	*iconGrid
}

func NewDungeonRewardGrid() (*DungeonRewardGrid, error) {
	// Reserve caption space for the widest label so reward icons align
	// across rows.
	zone := 0
	for _, s := range rewardLabels {
		zone = max(zone, icons.TextWidth(s))
	}
	captionedIcon := func(i int) (image.Image, error) {
		img, err := icons.DungeonRewardIcon(i)
		if err != nil {
			return nil, err
		}
		return icons.Captioned(rewardLabels[i], zone, img), nil
	}

	// Captioned images are uniformly sized; cells size to fit them.
	g, err := newIconGrid(icons.DungeonRewardCount(), rewardGridColumns, 0, 0, captionedIcon)
	if err != nil {
		return nil, err
	}
	grid := &DungeonRewardGrid{iconGrid: g}

	for i, label := range g.labels {
		reward, err := icons.DungeonRewardAt(i)
		if err != nil {
			return nil, err
		}
		Bind(label, "<Button-1>", Command(func() {
			reward.ToggleAcquired()
			grid.Refresh(i)
		}))
		Bind(label, "<Button-3>", Command(func() {
			reward.CycleKind()
			grid.Refresh(i)
		}))
	}

	return grid, nil
}
