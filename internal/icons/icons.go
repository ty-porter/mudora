package icons

import (
	"embed"
	"strings"
)

//go:embed assets/*.png
var assets embed.FS

// itemFile maps an item name (as it appears in items.txt / rom.Inspect) to its
// sprite file under assets/, without the .png suffix. Dungeon-specific keys,
// maps and compasses are resolved by prefix in fileFor instead of listed here.
var itemFile = map[string]string{
	// Swords and shields.
	"Progressive Sword":  "sword",
	"Master Sword":       "sword",
	"Progressive Shield": "shield",

	// Mail.
	"Progressive Mail":  "mail-green",
	"Progressive Armor": "mail-green",

	// Bow and arrows.
	"Bow":                "bow",
	"Progressive Bow":    "bow",
	"Arrow":              "arrow-1",
	"Arrows (10)":        "arrow-10",
	"Arrow Upgrade (5)":  "arrow-5",
	"Arrow Upgrade (10)": "arrow-10",
	"Arrow Upgrade (70)": "arrow-10",

	// Boomerangs.
	"Boomerang":       "boomerang-blue",
	"Boomerang (Red)": "boomerang-red",

	// Rods, medallions, tools, equipment.
	"Hookshot":          "hookshot",
	"Fire Rod":          "fire-rod",
	"Ice Rod":           "ice-rod",
	"Bombos":            "bombos",
	"Ether":             "ether",
	"Quake":             "quake",
	"Lamp":              "lamp",
	"Hammer":            "hammer",
	"Shovel":            "shovel",
	"Powder":            "powder",
	"Mushroom":          "mushroom",
	"Ocarina":           "flute",
	"Flute":             "flute",
	"Bug Net":           "bugnet",
	"Book of Mudora":    "book",
	"Cane of Somaria":   "somaria",
	"Cane of Byrna":     "byrna",
	"Cape":              "cape",
	"Magic Mirror":      "magic-mirror",
	"Moon Pearl":        "moonpearl",
	"Pegasus Boots":     "boots",
	"Flippers":          "flippers",
	"Progressive Glove": "glove",

	// Bombs.
	"Bomb":              "bomb-1",
	"Bombs (3)":         "bomb-3",
	"Bombs (10)":        "bomb-10",
	"Bomb Upgrade (5)":  "bomb-5",
	"Bomb Upgrade (10)": "bomb-10",
	"Bomb Upgrade (50)": "bomb-10",

	// Bottles and standalone potions.
	"Bottle":                "bottle-empty",
	"Bottle (Red Potion)":   "bottle-red",
	"Bottle (Green Potion)": "bottle-green",
	"Bottle (Blue Potion)":  "bottle-blue",
	"Bottle (Bee)":          "bottle-bee",
	"Bottle (Super bee)":    "bottle-bee",
	"Bottle (Faerie)":       "bottle-faerie",
	"Red Potion":            "bottle-red",
	"Green Potion":          "bottle-green",
	"Blue Potion":           "bottle-blue",

	// Magic.
	"Small Magic":   "magic-small",
	"Half Magic":    "magic-large",
	"Quarter Magic": "magic-large",

	// Hearts.
	"Heart":           "heart",
	"Heart Container": "heart-container",
	"Piece of Heart":  "heart-piece",

	// Rupees.
	"Rupee (Green)": "rupee-green-1",
	"Rupee (Blue)":  "rupee-blue-5",
	"Rupee (Red)":   "rupee-red-20",
	"Rupees (20)":   "rupee-green-20",
	"Rupees (50)":   "rupee-green-50",
	"Rupees (100)":  "rupee-green-100",
	"Rupees (300)":  "rupees-green-300",

	// Pendants and crystals.
	"Pendant of Courage": "pendant-green",
	"Pendant of Power":   "pendant-red",
	"Pendant of Wisdom":  "pendant-blue",
	"Crystal":            "crystal",
}

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

func fileFor(item string) (string, bool) {
	if f, ok := itemFile[item]; ok {
		return f, true
	}
	switch {
	case strings.HasPrefix(item, "Big Key"):
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
