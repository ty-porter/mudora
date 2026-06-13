package icons

import (
	"fmt"
	"image"
)

// Tracker items. These are pointers so that mutating State (e.g. from the
// auto-tracker) is seen by everything holding the item, including the items
// table below.
var (
	Sword = &GridItem{
		boundingBoxes: []Box{
			{X: 81, Y: 10, W: 6, H: 13},
			{X: 97, Y: 8, W: 7, H: 16},
			{X: 129, Y: 8, W: 7, H: 16},
			{X: 145, Y: 8, W: 7, H: 16},
		},
	}
	Shield = &GridItem{
		boundingBoxes: []Box{
			{X: 160, Y: 11, W: 8, H: 10},
			{X: 179, Y: 12, W: 11, H: 12},
			{X: 201, Y: 8, W: 14, H: 16},
		},
	}
	Bow = &GridItem{
		boundingBoxes: []Box{
			{X: 84, Y: 33, W: 15, H: 15},
			{X: 196, Y: 33, W: 15, H: 15},
		},
	}
	Boomerang = &GridItem{
		boundingBoxes: []Box{
			{X: 24, Y: 98, W: 8, H: 12},
			{X: 40, Y: 98, W: 8, H: 12},
		},
	}
	Hookshot = &GridItem{
		boundingBoxes: []Box{
			{X: 56, Y: 96, W: 7, H: 16},
		},
	}
	Bomb = &GridItem{
		boundingBoxes: []Box{
			{X: 50, Y: 57, W: 13, H: 14},
		},
	}
	Powder = &GridItem{
		boundingBoxes: []Box{
			{X: 96, Y: 97, W: 16, H: 15},
		},
	}
	Mushroom = &GridItem{
		boundingBoxes: []Box{
			{X: 72, Y: 96, W: 16, H: 16},
		},
	}
	FireRod = &GridItem{
		boundingBoxes: []Box{
			{X: 120, Y: 96, W: 8, H: 16},
		},
	}
	IceRod = &GridItem{
		boundingBoxes: []Box{
			{X: 136, Y: 96, W: 8, H: 16},
		},
	}
	Bombos = &GridItem{
		boundingBoxes: []Box{
			{X: 152, Y: 96, W: 16, H: 16},
		},
	}
	Ether = &GridItem{
		boundingBoxes: []Box{
			{X: 176, Y: 96, W: 16, H: 16},
		},
	}
	Quake = &GridItem{
		boundingBoxes: []Box{
			{X: 200, Y: 96, W: 16, H: 16},
		},
	}
	Lamp = &GridItem{
		boundingBoxes: []Box{
			{X: 226, Y: 96, W: 12, H: 16},
		},
	}
	Shovel = &GridItem{
		boundingBoxes: []Box{
			{X: 61, Y: 120, W: 8, H: 16},
		},
	}
	Hammer = &GridItem{
		boundingBoxes: []Box{
			{X: 248, Y: 97, W: 8, H: 15},
		},
	}
	// Flute has an on/off state
	Flute = &GridItem{
		boundingBoxes: []Box{
			{X: 78, Y: 121, W: 14, H: 14},
			{X: 78, Y: 121, W: 14, H: 14},
		},
	}
	BugNet = &GridItem{
		boundingBoxes: []Box{
			{X: 102, Y: 120, W: 13, H: 16},
		},
	}
	Book = &GridItem{
		boundingBoxes: []Box{
			{X: 126, Y: 121, W: 13, H: 15},
		},
	}
	MoonPearl = &GridItem{
		boundingBoxes: []Box{
			{X: 210, Y: 170, W: 12, H: 12},
		},
	}
	Somaria = &GridItem{
		boundingBoxes: []Box{
			{X: 149, Y: 120, W: 8, H: 16},
		},
	}
	Byrna = &GridItem{
		boundingBoxes: []Box{
			{X: 165, Y: 120, W: 8, H: 16},
		},
	}
	Cape = &GridItem{
		boundingBoxes: []Box{
			{X: 181, Y: 120, W: 16, H: 16},
		},
	}
	Mirror = &GridItem{
		boundingBoxes: []Box{
			{X: 206, Y: 121, W: 14, H: 15},
		},
	}
	Mail = &GridItem{
		boundingBoxes: []Box{
			{X: 40, Y: 170, W: 16, H: 14},
			{X: 64, Y: 170, W: 16, H: 14},
			{X: 88, Y: 170, W: 16, H: 14},
		},
	}
	Boots = &GridItem{
		boundingBoxes: []Box{
			{X: 113, Y: 169, W: 14, H: 15},
		},
	}
	Glove = &GridItem{
		boundingBoxes: []Box{
			{X: 137, Y: 168, W: 14, H: 16},
			{X: 161, Y: 168, W: 14, H: 16},
		},
	}
	Flippers = &GridItem{
		boundingBoxes: []Box{
			{X: 185, Y: 168, W: 15, H: 16},
		},
	}
	Bottle = &GridItem{
		boundingBoxes: []Box{
			{X: 37, Y: 144, W: 14, H: 16},
		},
	}
	HalfMagic = &GridItem{
		boundingBoxes: []Box{
			{X: 136, Y: 273, W: 8, H: 15},
		},
	}
)

var items = []*GridItem{
	Bow, Boomerang, Hookshot, Bomb, Powder, Mushroom,
	FireRod, IceRod, Bombos, Ether, Quake, Shovel,
	Lamp, Hammer, Flute, BugNet, Book, HalfMagic,
	Bottle, Somaria, Byrna, Cape, Mirror, MoonPearl,
	Sword, Shield, Mail, Glove, Boots, Flippers,
}

// Count is the number of defined item icons.
func Count() int { return len(items) }

// ItemAt returns the grid item shown in cell id.
func ItemAt(id int) (*GridItem, error) {
	if id < 0 || id >= len(items) {
		return nil, fmt.Errorf("icons: no item %d", id)
	}
	return items[id], nil
}

// Item returns the icon for an item id in its current state. Disabled items
// (state 0) are grayed out.
func Item(id int) (image.Image, error) {
	if id < 0 || id >= len(items) {
		return nil, fmt.Errorf("icons: no item %d", id)
	}
	return items[id].Icon()
}
