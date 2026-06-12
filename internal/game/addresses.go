package game

const (
	WRAMBase uint32 = 0xF50000 // SNES WRAM $7E0000 mirrored here
	SRAMBase uint32 = 0xE00000 // cartridge SRAM
)

// In ALTTP the live save/inventory block sits in WRAM at $7EF000-$7EF4FF
// (a copy of the SRAM save slot). All item/dungeon state the tracker cares
// about lives in this one region, so a single read covers it.
const (
	SaveDataAddr uint32 = WRAMBase + 0xF000
	SaveDataSize uint32 = 0x500
)

// Offsets into the save data block ($7EF000 + offset).
// Reference: https://alttp-wiki.net/index.php/SRAM_Map
const (
	OffBow       = 0x340
	OffBoomerang = 0x341
	OffHookshot  = 0x342
	OffBombs     = 0x343
	OffPowder    = 0x344 // mushroom/powder
	OffFireRod   = 0x345
	OffIceRod    = 0x346
	OffBombos    = 0x347
	OffEther     = 0x348
	OffQuake     = 0x349
	OffLamp      = 0x34A
	OffHammer    = 0x34B
	OffFlute     = 0x34C // shovel/flute
	OffBugNet    = 0x34D
	OffBook      = 0x34E
	OffBottles   = 0x34F
	OffSomaria   = 0x350
	OffByrna     = 0x351
	OffCape      = 0x352
	OffMirror    = 0x353
	OffGloves    = 0x354
	OffBoots     = 0x355
	OffFlippers  = 0x356
	OffMoonPearl = 0x357
	OffSword     = 0x359
	OffShield    = 0x35A
	OffArmor     = 0x35B
	OffBottle1   = 0x35C
	OffBottle2   = 0x35D
	OffBottle3   = 0x35E
	OffBottle4   = 0x35F
	OffPendants  = 0x374
	OffCrystals  = 0x37A
	OffMagicUse  = 0x37B // 0=normal, 1=1/2, 2=1/4
	// TODO: compasses/big keys/maps (0x364-0x369), small key counts,
	// progressive item flags (0x38C-0x38E), ALTTPR-specific extensions, etc.
)

// Game module ($7E0010) — useful for knowing whether we're in-game at all
// (title screen / file select / triforce room reads garbage otherwise).
const (
	GameModeAddr uint32 = WRAMBase + 0x0010
)

// Game module values worth knowing about.
const (
	ModuleOverworld byte = 0x09
	ModuleDungeon   byte = 0x07
	// TODO: fill in the rest (title screen 0x00, file select 0x01,
	// triforce/credits 0x19/0x1A, ...).
)
