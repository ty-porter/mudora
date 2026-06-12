package ui

import (
	"github.com/ty-porter/mudora/internal/game"
	"github.com/ty-porter/mudora/internal/ui/icons"
)

// applyState maps a parsed game.State onto the item icons' display tiers.
// Progressive items resolve to a tier index (0 = not obtained, grayed out);
// the rest are simple obtained/not toggles. Items the save block doesn't track
// as inventory (e.g. bombs) are left untouched so manual edits survive.
func applyState(s game.State) {
	icons.Bow.State = bowTier(s.Bow)
	icons.Boomerang.State = clampTier(int(s.Boomerang), 2)
	icons.Hookshot.State = boolTier(s.Hookshot)
	icons.Powder.State = boolTier(s.Powder == 2)
	icons.Mushroom.State = boolTier(s.Powder == 1)
	icons.FireRod.State = boolTier(s.FireRod)
	icons.IceRod.State = boolTier(s.IceRod)
	icons.Bombos.State = boolTier(s.Bombos)
	icons.Ether.State = boolTier(s.Ether)
	icons.Quake.State = boolTier(s.Quake)
	icons.Shovel.State = boolTier(s.Flute == 1)
	icons.Lamp.State = boolTier(s.Lamp)
	icons.Hammer.State = boolTier(s.Hammer)
	icons.Flute.State = boolTier(s.Flute >= 2)
	icons.BugNet.State = boolTier(s.BugNet)
	icons.Book.State = boolTier(s.Book)
	icons.HalfMagic.State = boolTier(s.MagicUse > 0)
	icons.Bottle.State = boolTier(hasBottle(s.Bottles))
	icons.Somaria.State = boolTier(s.Somaria)
	icons.Byrna.State = boolTier(s.Byrna)
	icons.Cape.State = boolTier(s.Cape)
	icons.Mirror.State = boolTier(s.Mirror)
	icons.MoonPearl.State = boolTier(s.MoonPearl)
	icons.Sword.State = clampTier(int(s.Sword), 4)
	icons.Shield.State = clampTier(int(s.Shield), 3)
	// Green mail (Armor 0) is the starting armor, so it always shows in color;
	// blue/red shift up one tier.
	icons.Mail.State = clampTier(int(s.Armor)+1, 3)
	icons.Glove.State = clampTier(int(s.Gloves), 2)
	icons.Boots.State = boolTier(s.Boots)
	icons.Flippers.State = boolTier(s.Flippers)
}

// bowTier maps the SRAM bow byte to the Bow icon's two tiers: 0 none, tier 1
// wooden bow, tier 2 silver (byte >= 3, with or without arrows).
func bowTier(v byte) int {
	switch {
	case v == 0:
		return 0
	case v >= 3:
		return 2
	default:
		return 1
	}
}

func boolTier(b bool) int {
	if b {
		return 1
	}
	return 0
}

func clampTier(v, hi int) int {
	if v < 0 {
		return 0
	}
	if v > hi {
		return hi
	}
	return v
}

func hasBottle(bottles [4]byte) bool {
	for _, b := range bottles {
		if b != 0 {
			return true
		}
	}
	return false
}
