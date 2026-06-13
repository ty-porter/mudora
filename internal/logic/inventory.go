// Package logic models ALttP item-placement accessibility: given the player's
// current items and progress, which item-check locations can be reached.
//
// The rules here target the randomizer's glitchless "advanced item placement"
// ruleset (the same logic EmoTracker's ALttP pack encodes). Rules are written
// against Progress, a thin view over a parsed game.State plus the few world
// flags that aren't inventory items.
package logic

import (
	"math/bits"

	"github.com/ty-porter/mudora/internal/game"
)

// Inv is the access-logic view over a parsed game.State. It adds the derived
// predicates ALttP location rules are written in terms of ("can light a fire",
// "can lift heavy rocks", ...), keeping the rules readable and the SRAM
// encoding (progressive byte values, bitmasks) in one place.
type Inv struct {
	game.State
}

// For wraps a parsed game state for use by the access rules.
func For(s game.State) Inv { return Inv{State: s} }

// --- Single items -----------------------------------------------------------

func (i Inv) Sword() bool        { return i.State.Sword > 0 }
func (i Inv) MasterSword() bool  { return i.State.Sword >= 2 }
func (i Inv) Bow() bool          { return i.State.Bow > 0 }
func (i Inv) SilverArrows() bool { return i.State.Bow >= 3 }
func (i Inv) Boomerang() bool    { return i.State.Boomerang > 0 }
func (i Inv) RedBoomerang() bool { return i.State.Boomerang >= 2 }
func (i Inv) Hookshot() bool     { return i.State.Hookshot }
func (i Inv) Bombs() bool        { return i.State.Bombs }
func (i Inv) Mushroom() bool     { return i.State.Powder == 1 }
func (i Inv) Powder() bool       { return i.State.Powder == 2 }
func (i Inv) FireRod() bool      { return i.State.FireRod }
func (i Inv) IceRod() bool       { return i.State.IceRod }
func (i Inv) Bombos() bool       { return i.State.Bombos }
func (i Inv) Ether() bool        { return i.State.Ether }
func (i Inv) Quake() bool        { return i.State.Quake }
func (i Inv) Lamp() bool         { return i.State.Lamp }
func (i Inv) Hammer() bool       { return i.State.Hammer }
func (i Inv) Shovel() bool       { return i.State.Flute == 1 }
func (i Inv) Flute() bool        { return i.State.Flute >= 2 }
func (i Inv) Net() bool          { return i.State.BugNet }
func (i Inv) Book() bool         { return i.State.Book }
func (i Inv) Somaria() bool      { return i.State.Somaria }
func (i Inv) Byrna() bool        { return i.State.Byrna }
func (i Inv) Cape() bool         { return i.State.Cape }
func (i Inv) Mirror() bool       { return i.State.Mirror }
func (i Inv) Gloves() bool       { return i.State.Gloves > 0 }  // Power Gloves or better
func (i Inv) Mitts() bool        { return i.State.Gloves >= 2 } // Titan's Mitts
func (i Inv) Boots() bool        { return i.State.Boots }
func (i Inv) Flippers() bool     { return i.State.Flippers }
func (i Inv) MoonPearl() bool    { return i.State.MoonPearl }

// Agahnim reports whether Agahnim 1 has been beaten (opens the Dark World).
// It's a game-progress flag rather than an item, parsed from SRAM into State.
func (i Inv) Agahnim() bool { return i.State.AgahnimDefeated }

// Bottle reports whether the player holds at least one bottle (in any state).
func (i Inv) Bottle() bool { return i.BottleCount() > 0 }

// BottleCount is how many bottles the player holds (0-4).
func (i Inv) BottleCount() int {
	n := 0
	for _, c := range i.State.Bottles {
		if c != 0 {
			n++
		}
	}
	return n
}

// SwordLevel and GloveLevel expose the progressive tiers the access rules
// threshold on (e.g. sword2 = Master Sword).
func (i Inv) SwordLevel() int { return int(i.State.Sword) }
func (i Inv) GloveLevel() int { return int(i.State.Gloves) }

// PrizeCount is the total pendants + crystals collected (the rules' "prize").
func (i Inv) PrizeCount() int { return i.PendantCount() + i.CrystalCount() }

// HalfMagic / QuarterMagic report the magic-meter upgrade tier.
func (i Inv) HalfMagic() bool    { return i.State.MagicUse >= 1 }
func (i Inv) QuarterMagic() bool { return i.State.MagicUse >= 2 }

// PendantCount and CrystalCount report how many prizes have been collected.
// TODO: per-prize helpers (e.g. green pendant for Sahasrahla) need the bit
// layout of State.Pendants/State.Crystals, which package game doesn't yet
// document.
func (i Inv) PendantCount() int { return bits.OnesCount8(i.State.Pendants) }
func (i Inv) CrystalCount() int { return bits.OnesCount8(i.State.Crystals) }

// --- Derived capabilities ---------------------------------------------------

// CanLiftRocks reports whether the player can lift light rocks/bushes.
func (i Inv) CanLiftRocks() bool { return i.Gloves() }

// CanLiftHeavyRocks reports whether the player can lift dark rocks.
func (i Inv) CanLiftHeavyRocks() bool { return i.Mitts() }

// CanLightFires reports whether the player can light torches / dark rooms.
func (i Inv) CanLightFires() bool { return i.FireRod() || i.Lamp() }

// CanMeleeAttack reports whether the player has a basic melee option.
func (i Inv) CanMeleeAttack() bool { return i.Sword() || i.Hammer() }

// CanShootArrows reports whether the player can fire arrows.
func (i Inv) CanShootArrows() bool { return i.Bow() }

// CanMeltThings reports whether the player can melt the frozen barriers /
// freezor-style obstacles (fire rod, or bombos with a way to fight).
func (i Inv) CanMeltThings() bool { return i.FireRod() || (i.Bombos() && i.Sword()) }
