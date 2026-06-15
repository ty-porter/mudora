package rom

import (
	_ "embed"
)

//go:embed locations.txt
var locationData string

var Locations = parseAddresses("locations.txt", locationData)

func LocationAt(addr uint32) (string, bool) {
	name, ok := Locations[addr]
	return name, ok
}
