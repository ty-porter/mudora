package logic

// Settings are the randomizer options the pack's access rules branch on. The
// zero value, via DefaultSettings, describes a standard glitchless game: open
// world, no entrance/overworld/door shuffle, 7-crystal requirements.
//
// The codetracker pack also branches on the player's physical position via
// "toggle" provider codes (ow_dark_witch, hammer_pegs, ...). This interpreter
// has no position autotracking, so those toggles always resolve to 0 and
// reachability is computed purely from items along the rules' item paths.
type Settings struct {
	Inverted      bool // world_state_inverted (else world_state_open)
	Swordless     bool
	OWGlitches    bool // glitch_mode_ow
	MajorGlitches bool // glitch_mode_major
	Entrance      bool // any entrance shuffle (else entrance_shuffle_off)
	Overworld     bool // overworld shuffle (else ow_shuffle_off)
	Doors         bool // door shuffle (else door_shuffle_off)
	GTCrystals    int  // crystals required to open Ganon's Tower
	GanonCrystals int  // crystals required to damage Ganon
}

// DefaultSettings is a standard open, glitchless, no-shuffle 7/7 game.
func DefaultSettings() Settings {
	return Settings{GTCrystals: 7, GanonCrystals: 7}
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

// settingCount resolves a settings/mode pseudo-code to a provider count (1 or
// 0), reporting whether the code was a recognized setting.
func (s Settings) settingCount(code string) (int, bool) {
	switch code {
	case "world_state_open", "world_state_standard":
		return b2i(!s.Inverted), true
	case "world_state_inverted":
		return b2i(s.Inverted), true
	case "swordless":
		return b2i(s.Swordless), true
	case "glitch_mode_ow":
		return b2i(s.OWGlitches), true
	case "glitch_mode_major":
		return b2i(s.MajorGlitches), true
	case "glitch_mode_none":
		return b2i(!s.OWGlitches && !s.MajorGlitches), true
	case "entrance_shuffle_off":
		return b2i(!s.Entrance), true
	case "entrance_shuffle_singleentryvanilla":
		// Vanilla single-entry caves: true unless full entrance shuffle.
		return b2i(!s.Entrance), true
	case "entrance_shuffle_on", "entrance_shuffle_all", "entrance_shuffle_crossed",
		"entrance_shuffle_dungeon":
		return b2i(s.Entrance), true
	case "ow_shuffle_off":
		return b2i(!s.Overworld), true
	case "ow_shuffle_on":
		return b2i(s.Overworld), true
	case "door_shuffle_off":
		return b2i(!s.Doors), true
	case "door_shuffle_sameset", "door_shuffle_on", "door_shuffle_crossed":
		return b2i(s.Doors), true
	}
	return 0, false
}

// itemCount resolves an item provider code to the player's count of it. The
// many EmoTracker aliases for each item are folded onto the same value.
func itemCount(code string, inv Inv) (int, bool) {
	switch code {
	case "bow":
		return b2i(inv.Bow()), true
	case "silvers", "silverarrows":
		return b2i(inv.SilverArrows()), true
	case "bombs", "bomb":
		return b2i(inv.Bombs()), true
	case "boomerang", "blue_boomerang", "bluemarang":
		return b2i(inv.Boomerang()), true
	case "red_boomerang", "redmarang", "magicboomerang", "magic_boomerang", "redboomerang":
		return b2i(inv.RedBoomerang()), true
	case "hookshot", "hs":
		return b2i(inv.Hookshot()), true
	case "powder":
		return b2i(inv.Powder()), true
	case "mushroom":
		return b2i(inv.Mushroom()), true
	case "firerod":
		return b2i(inv.FireRod()), true
	case "icerod":
		return b2i(inv.IceRod()), true
	case "bombos":
		return b2i(inv.Bombos()), true
	case "ether":
		return b2i(inv.Ether()), true
	case "quake":
		return b2i(inv.Quake()), true
	case "lamp":
		return b2i(inv.Lamp()), true
	case "hammer":
		return b2i(inv.Hammer()), true
	case "shovel":
		return b2i(inv.Shovel()), true
	case "flute", "activated_flute", "fluteactivated":
		return b2i(inv.Flute()), true
	case "net", "bugnet":
		return b2i(inv.Net()), true
	case "book":
		return b2i(inv.Book()), true
	case "somaria":
		return b2i(inv.Somaria()), true
	case "byrna":
		return b2i(inv.Byrna()), true
	case "cape":
		return b2i(inv.Cape()), true
	case "mirror":
		return b2i(inv.Mirror()), true
	case "boots", "pegasus":
		return b2i(inv.Boots()), true
	case "flippers":
		return b2i(inv.Flippers()), true
	case "moonpearl", "pearl":
		return b2i(inv.MoonPearl()), true
	case "glove", "lift1", "powerglove":
		return b2i(inv.Gloves()), true
	case "mitt", "lift2", "titansmitt", "titansmitts":
		return b2i(inv.Mitts()), true
	case "sword", "sword1":
		return b2i(inv.SwordLevel() >= 1), true
	case "sword2", "mastersword":
		return b2i(inv.SwordLevel() >= 2), true
	case "sword3":
		return b2i(inv.SwordLevel() >= 3), true
	case "sword4":
		return b2i(inv.SwordLevel() >= 4), true
	case "halfmagic":
		return b2i(inv.HalfMagic()), true
	case "quartermagic":
		return b2i(inv.QuarterMagic()), true
	case "bottle":
		return inv.BottleCount(), true
	case "aga", "agahnim", "aga1":
		return b2i(inv.Agahnim()), true
	case "crystal", "crystals":
		return inv.CrystalCount(), true
	case "pendant", "pendants":
		return inv.PendantCount(), true
	case "prize", "prizes":
		return inv.PrizeCount(), true
	}
	return 0, false
}
