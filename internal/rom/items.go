package rom

import (
	_ "embed"
)

//go:embed items.txt
var itemData string

var Items = parseAddresses("items.txt", itemData)

func ItemAt(addr uint32) (string, bool) {
	name, ok := Items[addr]
	return name, ok
}
