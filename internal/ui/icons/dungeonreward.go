package icons

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

// Reward kinds, in the order right-click cycles through them.
const (
	Crystal = iota
	RedCrystal
	GreenPendant
	BluePendant
	RedPendant
	rewardKindCount
)

// rewardBoxes maps reward kind to its sprite. The sheet has no red-crystal
// sprite; RedCrystal reuses the crystal box and is marked with a red plus
// when drawn.
var rewardBoxes = [rewardKindCount]Box{
	Crystal:      {X: 159, Y: 193, W: 10, H: 14},
	RedCrystal:   {X: 159, Y: 193, W: 10, H: 14},
	GreenPendant: {X: 84, Y: 193, W: 16, H: 15},
	BluePendant:  {X: 108, Y: 193, W: 16, H: 15},
	RedPendant:   {X: 132, Y: 193, W: 16, H: 15},
}

// DungeonReward is one cell of the reward grid: which reward the dungeon
// holds and whether it has been acquired. Unacquired rewards draw grayed
// out.
type DungeonReward struct {
	Kind     int
	Acquired bool
}

// Icon returns the sprite for the dungeon's current reward, upscaled for
// display.
func (d *DungeonReward) Icon() (image.Image, error) {
	if d.Kind < 0 || d.Kind >= rewardKindCount {
		return nil, fmt.Errorf("icons: reward kind %d out of range", d.Kind)
	}
	s, err := loadSheet()
	if err != nil {
		return nil, err
	}
	b := rewardBoxes[d.Kind]
	img := image.Image(s.SubImage(image.Rect(b.X, b.Y, b.X+b.W, b.Y+b.H)))
	if !d.Acquired {
		img = grayscale(img)
	}
	// Marked after graying so the plus stays visible on unacquired rewards.
	if d.Kind == RedCrystal {
		img = markPlus(img)
	}
	return upscale(img), nil
}

// CycleKind advances the reward to the next kind, wrapping around.
func (d *DungeonReward) CycleKind() {
	d.Kind = (d.Kind + 1) % rewardKindCount
}

// ToggleAcquired flips whether the reward has been acquired.
func (d *DungeonReward) ToggleAcquired() {
	d.Acquired = !d.Acquired
}

// One tracker cell per reward-bearing dungeon.
var (
	EasternPalace    = &DungeonReward{}
	DesertPalace     = &DungeonReward{}
	TowerOfHera      = &DungeonReward{}
	PalaceOfDarkness = &DungeonReward{}
	SwampPalace      = &DungeonReward{}
	SkullWoods       = &DungeonReward{}
	ThievesTown      = &DungeonReward{}
	IcePalace        = &DungeonReward{}
	MiseryMire       = &DungeonReward{}
	TurtleRock       = &DungeonReward{}
)

// dungeonRewards defines the reward grid contents in display order.
var dungeonRewards = []*DungeonReward{
	EasternPalace,
	DesertPalace,
	TowerOfHera,
	PalaceOfDarkness,
	SwampPalace,
	SkullWoods,
	ThievesTown,
	IcePalace,
	MiseryMire,
	TurtleRock,
}

// DungeonRewardCount is the number of reward-bearing dungeons.
func DungeonRewardCount() int { return len(dungeonRewards) }

// DungeonRewardAt returns the dungeon reward shown in grid cell id.
func DungeonRewardAt(id int) (*DungeonReward, error) {
	if id < 0 || id >= len(dungeonRewards) {
		return nil, fmt.Errorf("icons: no dungeon %d", id)
	}
	return dungeonRewards[id], nil
}

// DungeonRewardIcon returns the reward icon for a dungeon id in its current
// state. Unacquired rewards are grayed out.
func DungeonRewardIcon(id int) (image.Image, error) {
	d, err := DungeonRewardAt(id)
	if err != nil {
		return nil, err
	}
	return d.Icon()
}

// markPlus draws a small light plus with a black border in the top-right
// corner; it marks the red-crystal variant the sheet lacks a sprite for.
func markPlus(src image.Image) image.Image {
	const size = 5 // odd, so the plus has a center pixel
	light := color.NRGBA{R: 230, G: 230, B: 230, A: 255}
	dark := color.NRGBA{A: 255}

	b := src.Bounds()
	out := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(out, out.Bounds(), src, b.Min, draw.Src)

	// Inset one pixel from the corner so the border isn't clipped.
	x0, y0 := b.Dx()-size-1, 1
	mid := size / 2
	var pts []image.Point
	for i := 0; i < size; i++ {
		pts = append(pts, image.Pt(x0+i, y0+mid), image.Pt(x0+mid, y0+i))
	}
	for _, p := range pts {
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if image.Pt(p.X+dx, p.Y+dy).In(out.Bounds()) {
					out.SetNRGBA(p.X+dx, p.Y+dy, dark)
				}
			}
		}
	}
	for _, p := range pts {
		out.SetNRGBA(p.X, p.Y, light)
	}
	return out
}
