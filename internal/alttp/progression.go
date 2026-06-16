package alttp

var Progression = map[string]bool{
	"Progressive Sword":  true,
	"Master Sword":       true,
	"Progressive Shield": true,
	"Progressive Mail":   true,
	"Progressive Armor":  true,

	"Bow":             true,
	"Progressive Bow": true,

	"Hookshot":          true,
	"Hammer":            true,
	"Fire Rod":          true,
	"Ice Rod":           true,
	"Bombos":            true,
	"Ether":             true,
	"Quake":             true,
	"Lamp":              true,
	"Flippers":          true,
	"Moon Pearl":        true,
	"Magic Mirror":      true,
	"Pegasus Boots":     true,
	"Progressive Glove": true,
	"Cane of Somaria":   true,
	"Cane of Byrna":     true,
	"Cape":              true,
	"Mushroom":          true,
	"Powder":            true,
	"Ocarina":           true,
	"Flute":             true,
	"Book of Mudora":    true,
	"Bug Net":           true,
	"Half Magic":        true,

	"Bottle":                true,
	"Bottle (Red Potion)":   true,
	"Bottle (Green Potion)": true,
	"Bottle (Blue Potion)":  true,
	"Bottle (Bee)":          true,
	"Bottle (Super bee)":    true,
	"Bottle (Faerie)":       true,

	"Pendant of Courage": true,
	"Pendant of Power":   true,
	"Pendant of Wisdom":  true,
	"Crystal":            true,
}

func IsProgression(item string) bool { return Progression[item] }
