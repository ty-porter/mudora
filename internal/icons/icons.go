// Package icons serves small item sprites, one embedded PNG per item, keyed by
// the item-value names that rom.Inspect produces. Items with no matching sprite
// (clocks, triforce pieces, RNG/programmable junk, ...) return ok == false.
package icons

import (
	"embed"
	"strings"
)

//go:embed assets/*.png
var assets embed.FS

// itemFile maps a rom.Inspect item-value name to its sprite file (the base name
// under assets/, without the .png suffix). Dungeon-specific keys, maps and
// compasses are resolved by prefix in fileFor rather than listed here.
var itemFile = map[string]string{
	// Swords and shields.
	"L1Sword":           "sword",
	"L1SwordAndShield":  "sword",
	"L2Sword":           "sword",
	"L3Sword":           "sword",
	"L4Sword":           "sword",
	"MasterSword":       "sword",
	"ProgressiveSword":  "sword",
	"BlueShield":        "shield",
	"RedShield":         "shield",
	"MirrorShield":      "shield",
	"ProgressiveShield": "shield",

	// Mail (only a green sprite is available).
	"BlueMail":         "mail-green",
	"RedMail":          "mail-green",
	"ProgressiveArmor": "mail-green",

	// Bow and arrows.
	"Bow":                     "bow",
	"BowAndArrows":            "bow",
	"BowAndSilverArrows":      "bow",
	"ProgressiveBow":          "bow",
	"ProgressiveBowAlternate": "bow",
	"Arrow":                   "arrow-1",
	"TenArrows":               "arrow-10",
	"ArrowUpgrade5":           "arrow-5",
	"ArrowUpgrade10":          "arrow-10",
	"ArrowUpgrade70":          "arrow-10",

	// Boomerangs.
	"Boomerang":    "boomerang-blue",
	"RedBoomerang": "boomerang-red",

	// Rods, medallions, tools, equipment.
	"Hookshot":         "hookshot",
	"FireRod":          "fire-rod",
	"IceRod":           "ice-rod",
	"Bombos":           "bombos",
	"Ether":            "ether",
	"Quake":            "quake",
	"Lamp":             "lamp",
	"Hammer":           "hammer",
	"Shovel":           "shovel",
	"Powder":           "powder",
	"Mushroom":         "mushroom",
	"OcarinaActive":    "flute",
	"OcarinaInactive":  "flute",
	"BugCatchingNet":   "bugnet",
	"BookOfMudora":     "book",
	"CaneOfSomaria":    "somaria",
	"CaneOfByrna":      "byrna",
	"Cape":             "cape",
	"MagicMirror":      "magic-mirror",
	"MoonPearl":        "moonpearl",
	"PegasusBoots":     "boots",
	"Flippers":         "flippers",
	"PowerGlove":       "glove",
	"TitansMitt":       "glove",
	"ProgressiveGlove": "glove",

	// Bombs.
	"Bomb":          "bomb-1",
	"ThreeBombs":    "bomb-3",
	"TenBombs":      "bomb-10",
	"BombUpgrade5":  "bomb-5",
	"BombUpgrade10": "bomb-10",
	"BombUpgrade50": "bomb-10",

	// Bottles and standalone potions.
	"Bottle":                "bottle-empty",
	"BottleWithRedPotion":   "bottle-red",
	"BottleWithGreenPotion": "bottle-green",
	"BottleWithBluePotion":  "bottle-blue",
	"BottleWithBee":         "bottle-bee",
	"BottleWithGoldBee":     "bottle-bee",
	"BottleWithFairy":       "bottle-faerie",
	"RedPotion":             "bottle-red",
	"GreenPotion":           "bottle-green",
	"BluePotion":            "bottle-blue",

	// Magic.
	"SmallMagic":   "magic-small",
	"HalfMagic":    "magic-large",
	"QuarterMagic": "magic-large",

	// Hearts.
	"Heart":                     "heart",
	"HeartContainer":            "heart-container",
	"BossHeartContainer":        "heart-container",
	"HeartContainerNoAnimation": "heart-container",
	"PieceOfHeart":              "heart-piece",

	// Rupees.
	"OneRupee":           "rupee-green-1",
	"FiveRupees":         "rupee-blue-5",
	"TwentyRupees":       "rupee-red-20",
	"TwentyRupees2":      "rupee-green-20",
	"FiftyRupees":        "rupee-green-50",
	"OneHundredRupees":   "rupee-green-100",
	"ThreeHundredRupees": "rupees-green-300",

	// Pendants and crystals.
	"PendantOfCourage": "pendant-green",
	"PendantOfPower":   "pendant-red",
	"PendantOfWisdom":  "pendant-blue",
	"Crystal":          "crystal",
}

// PNG returns the sprite PNG for the given item-value name, and whether one
// exists. The bytes are suitable for tk9.0's NewPhoto.
func PNG(item string) ([]byte, bool) {
	file, ok := fileFor(item)
	if !ok {
		return nil, false
	}
	data, err := assets.ReadFile("assets/" + file + ".png")
	if err != nil {
		return nil, false
	}
	return data, true
}

// fileFor resolves an item name to a sprite base name. Dungeon-specific keys,
// big keys, maps and compasses all share a single sprite, matched by prefix.
func fileFor(item string) (string, bool) {
	if f, ok := itemFile[item]; ok {
		return f, true
	}
	switch {
	case strings.HasPrefix(item, "BigKey"):
		return "key-big", true
	case strings.HasPrefix(item, "Key"):
		return "key-small", true
	case strings.HasPrefix(item, "Compass"):
		return "compass", true
	case strings.HasPrefix(item, "Map"):
		return "map", true
	}
	return "", false
}
