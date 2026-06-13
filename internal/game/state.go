package game

type State struct {
	Bow       byte // 0 = none, 1 = Bow, 3 = Silvers
	Boomerang byte // 0 = none, 1 = blue, 2 = red
	Hookshot  bool
	Bombs     bool
	FireRod   bool
	IceRod    bool
	Bombos    bool
	Ether     bool
	Quake     bool
	Lamp      bool
	Hammer    bool
	Powder    byte // 0 = none, 1 = mushroom, 2 = magic powder
	Flute     byte // 0 = none, 1 = shovel, 2 = flute
	BugNet    bool
	Book      bool
	Somaria   bool
	Byrna     bool
	Cape      bool
	Mirror    bool
	Gloves    byte // 0 = none, 1 = Power Gloves, 2 = Titan's Mitts
	Boots     bool
	Flippers  bool
	MoonPearl bool
	Sword     byte    // 0 = none, 1 = Fighters, 2 = Master, 3 = Tempered, 4 = Golden
	Shield    byte    // 0 = none, 1 = Fighters, 2 = Red, 3 = Mirror
	Armor     byte    // 0 = green, 1 = blue, 2 = red
	Bottles   [4]byte // Count
	MagicUse  byte    // 0 = normal, 1 = 1/2, 2 = 1/4

	// Progression.
	Pendants        byte // bitmask
	Crystals        byte // bitmask
	AgahnimDefeated bool // Agahnim 1 beaten (opens the Dark World)
}

// InGame reports whether the given game module value represents actual
// gameplay (as opposed to title screen, file select, credits, ...), i.e.
// whether the save block in WRAM holds meaningful data.
func InGame(module byte) bool {
	return module >= 0x05 && module <= 0x1B
}
