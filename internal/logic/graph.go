package logic

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"

	"github.com/ty-porter/mudora/internal/game"
)

// The world model is transcribed verbatim from the EmoTracker codetracker
// pack: regions.json defines the region-accessibility graph, the location
// files define each item check and the region it sits in. They're embedded and
// interpreted at load time rather than hand-translated, so re-syncing with the
// pack is a matter of re-copying the JSON.
//
//go:embed pack/regions.json
var regionsJSON []byte

//go:embed pack/overworld.json
var overworldJSON []byte

//go:embed pack/underworld.json
var underworldJSON []byte

//go:embed pack/dungeons.json
var dungeonsJSON []byte

// locationFiles are the embedded check lists, interpreted in order.
var locationFiles = [][]byte{overworldJSON, underworldJSON, dungeonsJSON}

type rawRegion struct {
	Name        string      `json:"name"`
	AccessRules []string    `json:"access_rules"`
	Children    []rawRegion `json:"children"`
}

type rawLocation struct {
	Name         string        `json:"name"`
	Parent       string        `json:"parent"`
	Sections     []rawSection  `json:"sections"`
	Children     []rawLocation `json:"children"`
	MapLocations []rawMapLoc   `json:"map_locations"`
}

type rawSection struct {
	Name        string   `json:"name"`
	AccessRules []string `json:"access_rules"`
	ItemCount   int      `json:"item_count"`
}

// rawMapLoc is one placement of a location on a named map. A location may list
// several, gated by restrict_visibility_rules (e.g. which entrance-shuffle mode
// is active); the first whose restriction holds for the settings is used.
type rawMapLoc struct {
	Map      string   `json:"map"`
	X        int      `json:"x"`
	Y        int      `json:"y"`
	Restrict []string `json:"restrict_visibility_rules"`
}

// mapSourceSize is the pixel size of the pack's (and our) square overworld map
// images; map coordinates are in this space and normalized against it.
const mapSourceSize = 2007.0

// check is one obtainable item location: the region it lives in, the local
// rules (beyond reaching the region) needed to collect it, and where it sits
// on the overworld map (mapName empty if it isn't placed on one).
type check struct {
	Name    string
	Region  string
	Rules   []string
	MapName string
	X, Y    int
}

// World is the interpreted accessibility model: the region graph plus the flat
// list of checks, evaluated against an inventory under fixed Settings.
type World struct {
	settings Settings
	regions  map[string][]string
	checks   []check
}

// New builds the world model for the given randomizer settings.
func New(set Settings) (*World, error) {
	regions, err := loadRegions(regionsJSON)
	if err != nil {
		return nil, err
	}
	w := &World{settings: set, regions: regions}
	for _, data := range locationFiles {
		checks, err := loadChecks(data, set)
		if err != nil {
			return nil, err
		}
		w.checks = append(w.checks, checks...)
	}
	return w, nil
}

// loadRegions flattens the region tree (children included) into a name->rules
// map. Entries without access rules (colors, bare containers) are skipped.
func loadRegions(data []byte) (map[string][]string, error) {
	var raw []rawRegion
	if err := json.Unmarshal(stripJSONC(data), &raw); err != nil {
		return nil, fmt.Errorf("logic: parsing regions: %w", err)
	}
	out := make(map[string][]string)
	var walk func(r rawRegion)
	walk = func(r rawRegion) {
		if r.Name != "" && r.AccessRules != nil {
			out[r.Name] = r.AccessRules
		}
		for _, c := range r.Children {
			walk(c)
		}
	}
	for _, r := range raw {
		walk(r)
	}
	return out, nil
}

// loadChecks flattens a location file into checks. A section is a check when
// it holds an item (item_count > 0); single-section locations keep the
// location's name, multi-section ones are suffixed with the section name.
func loadChecks(data []byte, set Settings) ([]check, error) {
	var raw []rawLocation
	if err := json.Unmarshal(stripJSONC(data), &raw); err != nil {
		return nil, fmt.Errorf("logic: parsing locations: %w", err)
	}
	// Map placements are gated only by settings, so evaluate them with an
	// empty inventory; region refs shouldn't appear, but fail them closed.
	pick := evalContext{set: set, region: func(string) bool { return false }}
	var out []check
	var walk func(l rawLocation)
	walk = func(l rawLocation) {
		ml, placed := chooseMapLoc(l.MapLocations, pick)
		for _, s := range l.Sections {
			if s.ItemCount <= 0 {
				continue
			}
			name := l.Name
			if len(l.Sections) > 1 && s.Name != "" {
				name = l.Name + " - " + s.Name
			}
			c := check{Name: name, Region: l.Parent, Rules: s.AccessRules}
			if placed {
				c.MapName, c.X, c.Y = ml.Map, ml.X, ml.Y
			}
			out = append(out, c)
		}
		for _, c := range l.Children {
			walk(c)
		}
	}
	for _, l := range raw {
		walk(l)
	}
	return out, nil
}

// chooseMapLoc returns the first placement whose visibility restriction holds
// under the current settings (an empty restriction always holds).
func chooseMapLoc(locs []rawMapLoc, pick evalContext) (rawMapLoc, bool) {
	for _, ml := range locs {
		if len(ml.Restrict) == 0 || pick.evalRules(ml.Restrict) {
			return ml, true
		}
	}
	return rawMapLoc{}, false
}

// reachableRegions computes, as a monotonic fixpoint, the set of regions
// reachable with the given inventory. Regions reference each other (@Name), so
// reaching one can unlock others; relaxation repeats until nothing new opens.
func (w *World) reachableRegions(inv Inv) map[string]bool {
	acc := make(map[string]bool, len(w.regions))
	ctx := evalContext{inv: inv, set: w.settings, region: func(name string) bool { return acc[name] }}
	for changed := true; changed; {
		changed = false
		for name, rules := range w.regions {
			if acc[name] || len(rules) == 0 {
				continue
			}
			if ctx.evalRules(rules) {
				acc[name] = true
				changed = true
			}
		}
	}
	return acc
}

// context builds an evaluation context whose @refs resolve against the given
// reachable-region set.
func (w *World) context(inv Inv, acc map[string]bool) evalContext {
	return evalContext{inv: inv, set: w.settings, region: func(name string) bool { return acc[name] }}
}

// regionReachable reports whether a check's parent region is satisfied. A
// parent that isn't a defined region (a color group, or a dungeon "X Access"
// location referenced via @ rather than declared as a region) doesn't gate the
// check here; its section's own rules carry the real requirement.
func (w *World) regionReachable(region string, acc map[string]bool) bool {
	if _, ok := w.regions[region]; !ok {
		return true
	}
	return acc[region]
}

// AccessibleLocations returns the sorted names of every check obtainable with
// the given game state: its region is reachable and its local rules pass.
func (w *World) AccessibleLocations(s game.State) []string {
	inv := For(s)
	acc := w.reachableRegions(inv)
	ctx := w.context(inv, acc)
	var out []string
	for _, c := range w.checks {
		if w.regionReachable(c.Region, acc) && ctx.evalRules(c.Rules) {
			out = append(out, c.Name)
		}
	}
	sort.Strings(out)
	return out
}

// Marker is a location's placement on an overworld map, with its current
// logical accessibility. U and V are normalized [0,1] coordinates within the
// square map image.
type Marker struct {
	Name       string
	Map        string // e.g. "lightworld", "darkworld"
	U, V       float64
	Accessible bool
}

// Markers returns one marker per placed map position for the given game state,
// merging checks that share a position (accessible if any of them is). Checks
// not placed on a map are omitted.
func (w *World) Markers(s game.State) []Marker {
	inv := For(s)
	acc := w.reachableRegions(inv)
	ctx := w.context(inv, acc)

	type key struct {
		m    string
		x, y int
	}
	idx := map[key]int{}
	var out []Marker
	for _, c := range w.checks {
		if c.MapName == "" {
			continue
		}
		reach := w.regionReachable(c.Region, acc) && ctx.evalRules(c.Rules)
		k := key{c.MapName, c.X, c.Y}
		if i, ok := idx[k]; ok {
			out[i].Accessible = out[i].Accessible || reach
			continue
		}
		idx[k] = len(out)
		out = append(out, Marker{
			Name:       c.Name,
			Map:        c.MapName,
			U:          float64(c.X) / mapSourceSize,
			V:          float64(c.Y) / mapSourceSize,
			Accessible: reach,
		})
	}
	return out
}

// Accessible reports whether a single named check is obtainable. Unknown names
// report false.
func (w *World) Accessible(name string, s game.State) bool {
	inv := For(s)
	acc := w.reachableRegions(inv)
	ctx := w.context(inv, acc)
	for _, c := range w.checks {
		if c.Name == name {
			return w.regionReachable(c.Region, acc) && ctx.evalRules(c.Rules)
		}
	}
	return false
}

// UnknownRegions lists the check parent-regions that aren't defined in the
// region graph. Useful for spotting location files whose regions haven't been
// embedded yet; such checks can never be reachable.
func (w *World) UnknownRegions() []string {
	seen := map[string]bool{}
	var out []string
	for _, c := range w.checks {
		if _, ok := w.regions[c.Region]; !ok && !seen[c.Region] {
			seen[c.Region] = true
			out = append(out, c.Region)
		}
	}
	sort.Strings(out)
	return out
}

// trailingComma matches a comma that directly precedes a closing brace or
// bracket (a JSON5 trailing comma the standard decoder rejects).
var trailingComma = regexp.MustCompile(`,(\s*[}\]])`)

// stripJSONC turns the pack's JSON5-ish files (// and /* */ comments, trailing
// commas) into strict JSON the standard library can decode. Comment scanning
// is string-aware so // inside a string value is preserved.
func stripJSONC(src []byte) []byte {
	out := make([]byte, 0, len(src))
	inString, escaped := false, false
	for i := 0; i < len(src); i++ {
		c := src[i]
		if inString {
			out = append(out, c)
			switch {
			case escaped:
				escaped = false
			case c == '\\':
				escaped = true
			case c == '"':
				inString = false
			}
			continue
		}
		if c == '"' {
			inString = true
			out = append(out, c)
			continue
		}
		if c == '/' && i+1 < len(src) {
			switch src[i+1] {
			case '/':
				for i+1 < len(src) && src[i+1] != '\n' {
					i++
				}
				continue
			case '*':
				i += 2
				for i < len(src) && !(src[i] == '*' && i+1 < len(src) && src[i+1] == '/') {
					i++
				}
				i++ // skip the closing '/'
				continue
			}
		}
		out = append(out, c)
	}
	return trailingComma.ReplaceAll(out, []byte("$1"))
}
