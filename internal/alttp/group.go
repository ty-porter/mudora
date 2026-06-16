package alttp

import (
	"strings"

	"github.com/ty-porter/mudora/internal/rom"
)

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

func Filter(groups []Group, query string) []Group {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return groups
	}

	var out []Group
	for _, g := range groups {
		var matched []Placement
		for _, p := range g.Locations {
			if strings.Contains(strings.ToLower(p.Item), query) {
				matched = append(matched, p)
			}
		}
		if len(matched) > 0 {
			out = append(out, Group{Region: g.Region, Locations: matched})
		}
	}
	return out
}
