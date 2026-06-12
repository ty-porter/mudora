package game

import (
	"fmt"
)

func ParseSaveData(data []byte) (State, error) {
	if len(data) < int(SaveDataSize) {
		return State{}, fmt.Errorf("save data too short: got %d bytes, want %d", len(data), SaveDataSize)
	}

	s := State{
		Bow:       data[OffBow],
		Boomerang: data[OffBoomerang],
		Hookshot:  data[OffHookshot] != 0,
		FireRod:   data[OffFireRod] != 0,
		IceRod:    data[OffIceRod] != 0,
		Bombos:    data[OffBombos] != 0,
		Ether:     data[OffEther] != 0,
		Quake:     data[OffQuake] != 0,
		Lamp:      data[OffLamp] != 0,
		Hammer:    data[OffHammer] != 0,
		Powder:    data[OffPowder],
		Flute:     data[OffFlute],
		BugNet:    data[OffBugNet] != 0,
		Book:      data[OffBook] != 0,
		Somaria:   data[OffSomaria] != 0,
		Byrna:     data[OffByrna] != 0,
		Cape:      data[OffCape] != 0,
		Mirror:    data[OffMirror] != 0,
		Gloves:    data[OffGloves],
		Boots:     data[OffBoots] != 0,
		Flippers:  data[OffFlippers] != 0,
		MoonPearl: data[OffMoonPearl] != 0,
		Sword:     data[OffSword],
		Shield:    data[OffShield],
		Armor:     data[OffArmor],
		Bottles: [4]byte{
			data[OffBottle1],
			data[OffBottle2],
			data[OffBottle3],
			data[OffBottle4],
		},
		Pendants: data[OffPendants],
		Crystals: data[OffCrystals],
		MagicUse: data[OffMagicUse],
	}

	// TODO: ALTTPR stores some items progressively / in extended SRAM —
	// handle randomizer-specific encodings here (e.g. flute states,
	// mushroom vs powder both held, bow w/ silvers flags at 0x38E).

	return s, nil
}
