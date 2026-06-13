package logic

import (
	"testing"

	"github.com/ty-porter/mudora/internal/game"
)

func TestWorldLoads(t *testing.T) {
	w, err := New(DefaultSettings())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if len(w.regions) == 0 {
		t.Fatal("no regions loaded")
	}
	if len(w.checks) == 0 {
		t.Fatal("no checks loaded")
	}
	t.Logf("loaded %d regions, %d checks", len(w.regions), len(w.checks))
	if unknown := w.UnknownRegions(); len(unknown) > 0 {
		t.Logf("%d checks reference regions not in the graph: %v", len(unknown), unknown)
	}
}

func TestMarkers(t *testing.T) {
	w, err := New(DefaultSettings())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	markers := w.Markers(game.State{})
	if len(markers) == 0 {
		t.Fatal("no markers produced")
	}
	byMap := map[string]int{}
	reachable := 0
	for _, m := range markers {
		byMap[m.Map]++
		if m.U < 0 || m.U > 1 || m.V < 0 || m.V > 1 {
			t.Errorf("marker %q out of [0,1]: U=%v V=%v", m.Name, m.U, m.V)
		}
		if m.Accessible {
			reachable++
		}
	}
	t.Logf("markers: %d total %v; %d reachable with no items", len(markers), byMap, reachable)
	if byMap["lightworld"] == 0 || byMap["darkworld"] == 0 {
		t.Errorf("expected markers on both maps, got %v", byMap)
	}
}

func TestOverworldAccessibility(t *testing.T) {
	w, err := New(DefaultSettings())
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	state := func(mut func(s *game.State)) game.State {
		var s game.State
		if mut != nil {
			mut(&s)
		}
		return s
	}

	cases := []struct {
		name  string
		check string
		state game.State
		want  bool
	}{
		{
			// Reachable in open with no items (region anchored by @Light World,
			// rule is just world_state_open).
			name:  "Mushroom Spot, no items",
			check: "Mushroom Spot",
			state: state(nil),
			want:  true,
		},
		{
			// Pedestal needs the green+blue pendants or all 10 prizes; the
			// green-pendant path is unmodeled, so only prize:10 qualifies here.
			name:  "Pedestal with nothing",
			check: "Master Sword Pedestal",
			state: state(nil),
			want:  false,
		},
		{
			name:  "Pedestal with all 10 prizes",
			check: "Master Sword Pedestal",
			state: state(func(s *game.State) { s.Pendants = 0x07; s.Crystals = 0x7f }),
			want:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := w.Accessible(tc.check, tc.state); got != tc.want {
				t.Errorf("Accessible(%q) = %v, want %v", tc.check, got, tc.want)
			}
		})
	}

	// Sanity: with no items at all, a nonzero but non-total set of the
	// overworld is in logic.
	got := len(w.AccessibleLocations(state(nil)))
	t.Logf("accessible overworld checks with no items: %d / %d", got, len(w.checks))
	if got == 0 {
		t.Error("expected some checks reachable with no items in open mode")
	}
	if got == len(w.checks) {
		t.Error("expected some checks gated behind items, but all were reachable")
	}
}
