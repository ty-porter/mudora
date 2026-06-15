package rom

import "sort"

// Entry is the item placed at a known location in a ROM.
type Entry struct {
	Address  uint32
	Location string
	Item     string
}

func Inspect(data []byte) []Entry {
	entries := make([]Entry, 0, len(Locations))
	for addr, loc := range Locations {
		item := "UNKNOWN"
		if i := int(addr); i >= 0 && i < len(data) {
			if name, ok := ItemAt(uint32(data[i])); ok {
				item = name
			}
		}
		entries = append(entries, Entry{Address: addr, Location: loc, Item: item})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Address < entries[j].Address })
	return entries
}
