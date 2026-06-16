package alttp

import "github.com/ty-porter/mudora/internal/rom"

type Placement struct {
	Location string
	Item     string
}

type Group struct {
	Region    string
	Locations []Placement
}

func Grouped(entries []rom.Entry) []Group {
	itemByLoc := make(map[string]string, len(entries))
	for _, e := range entries {
		itemByLoc[e.Location] = e.Item
	}

	groups := make([]Group, 0, len(RegionOrder))
	for _, region := range RegionOrder {
		g := Group{Region: region}
		for _, loc := range Regions[region] {
			g.Locations = append(g.Locations, Placement{Location: loc, Item: itemByLoc[loc]})
		}
		groups = append(groups, g)
	}
	return groups
}
