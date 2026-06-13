package logic

import (
	"testing"

	"github.com/ty-porter/mudora/internal/game"
)

func TestRuleDSL(t *testing.T) {
	// A region resolver that only knows "Open Region" is reachable.
	regions := func(name string) bool { return name == "Open Region" }

	ctx := func(mut func(s *game.State), set Settings) evalContext {
		var s game.State
		if mut != nil {
			mut(&s)
		}
		return evalContext{inv: For(s), set: set, region: regions}
	}

	def := DefaultSettings()

	cases := []struct {
		name  string
		rules []string
		mut   func(s *game.State)
		set   Settings
		want  bool
	}{
		{"open world setting holds by default", []string{"world_state_open"}, nil, def, true},
		{"inverted setting is off by default", []string{"world_state_inverted"}, nil, def, false},
		{"AND: both items present", []string{"moonpearl,hammer"},
			func(s *game.State) { s.MoonPearl = true; s.Hammer = true }, def, true},
		{"AND: one item missing", []string{"moonpearl,hammer"},
			func(s *game.State) { s.MoonPearl = true }, def, false},
		{"OR: second rule passes", []string{"flippers", "boots"},
			func(s *game.State) { s.Boots = true }, def, true},
		{"count: prize:10 with all prizes", []string{"prize:10"},
			func(s *game.State) { s.Pendants = 0x07; s.Crystals = 0x7f }, def, true},
		{"count: prize:10 short", []string{"prize:10"},
			func(s *game.State) { s.Pendants = 0x07 }, def, false},
		{"sword2 needs master sword", []string{"sword2"},
			func(s *game.State) { s.Sword = 1 }, def, false},
		{"sword2 with master sword", []string{"sword2"},
			func(s *game.State) { s.Sword = 2 }, def, true},
		{"inspect tier never in logic", []string{"{book}"},
			func(s *game.State) { s.Book = true }, def, false},
		{"glitch tier off under glitchless", []string{"[flippers]"},
			func(s *game.State) { s.Flippers = true }, def, false},
		{"glitch tier on with OW glitches", []string{"[flippers]"},
			func(s *game.State) { s.Flippers = true }, Settings{OWGlitches: true, GTCrystals: 7, GanonCrystals: 7}, true},
		{"region ref resolves", []string{"@Open Region,hammer"},
			func(s *game.State) { s.Hammer = true }, def, true},
		{"region ref missing", []string{"@Closed Region"}, nil, def, false},
		{"func: barrier via cape", []string{"$canClearAgaTowerBarrier"},
			func(s *game.State) { s.Cape = true }, def, true},
		{"func: barrier needs master sword without cape", []string{"$canClearAgaTowerBarrier"},
			func(s *game.State) { s.Sword = 1 }, def, false},
		{"unknown position toggle resolves to 0", []string{"ow_dark_witch"}, nil, def, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ctx(tc.mut, tc.set).evalRules(tc.rules); got != tc.want {
				t.Errorf("evalRules(%v) = %v, want %v", tc.rules, got, tc.want)
			}
		})
	}
}
