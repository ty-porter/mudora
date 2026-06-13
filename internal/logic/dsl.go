package logic

import (
	"strconv"
	"strings"
)

// evalContext evaluates the pack's access-rule DSL against a fixed inventory
// and settings. Region cross-references (@Name) are resolved through region,
// which the graph supplies; this keeps the rule grammar testable in isolation.
type evalContext struct {
	inv    Inv
	set    Settings
	region func(name string) bool
}

// evalRules reports whether any single rule in the OR-list is satisfied.
func (c evalContext) evalRules(rules []string) bool {
	if len(rules) == 0 {
		return true // no rules listed = freely accessible
	}
	for _, r := range rules {
		if c.evalRule(r) {
			return true
		}
	}
	return false
}

// evalRule evaluates one comma-separated rule, all of whose tokens must hold.
func (c evalContext) evalRule(rule string) bool {
	rule = strings.TrimSpace(rule)
	if rule == "" {
		return true
	}
	for _, tok := range strings.Split(rule, ",") {
		if !c.evalToken(strings.TrimSpace(tok)) {
			return false
		}
	}
	return true
}

func (c evalContext) evalToken(tok string) bool {
	if tok == "" {
		return true
	}
	switch tok[0] {
	case '[': // glitch tier: only satisfiable when glitches are enabled
		inner := strings.TrimSuffix(tok[1:], "]")
		if c.set.OWGlitches || c.set.MajorGlitches {
			return c.evalToken(inner)
		}
		return false
	case '{': // inspect/visible tier: not "in logic"
		return false
	case '@': // another region's accessibility
		return c.region(tok[1:])
	case '$': // a ported logic_common helper
		name, arg := tok[1:], ""
		if i := strings.IndexByte(name, '|'); i >= 0 {
			name, arg = name[:i], name[i+1:]
		}
		return c.callFunc(name, arg) > 0
	}

	// Plain provider code, optionally "code:count".
	code, want := tok, 1
	if i := strings.LastIndexByte(tok, ':'); i >= 0 {
		if n, err := strconv.Atoi(tok[i+1:]); err == nil {
			code, want = tok[:i], n
		}
	}
	return c.count(code) >= want
}

// count resolves a bare provider code. Settings take precedence over items
// (disjoint namespaces); anything unrecognized is an unmodeled provider —
// chiefly the pack's position toggles — and resolves to 0.
func (c evalContext) count(code string) int {
	if n, ok := c.set.settingCount(code); ok {
		return n
	}
	if n, ok := itemCount(code, c.inv); ok {
		return n
	}
	return 0
}

// callFunc ports the logic_common.lua helpers the rules invoke. Unknown
// helpers return 0 (treated as unmet) so a missing port fails closed.
func (c evalContext) callFunc(name, arg string) int {
	switch name {
	case "canActivateTablets":
		if c.set.Swordless {
			return b2i(c.inv.Hammer())
		}
		return b2i(c.inv.SwordLevel() >= 2)
	case "canUseMedallions", "canRemoveCurtains":
		if c.set.Swordless {
			return 1
		}
		return b2i(c.inv.SwordLevel() >= 1)
	case "canClearAgaTowerBarrier":
		if c.inv.Cape() {
			return 1
		}
		if c.set.Swordless {
			return b2i(c.inv.Hammer())
		}
		return b2i(c.inv.SwordLevel() >= 2)
	case "gtCrystalCount":
		return b2i(c.inv.CrystalCount() >= c.set.GTCrystals || c.inv.PrizeCount() == 10)
	case "ganonCrystalCount":
		return b2i(c.inv.CrystalCount() >= c.set.GanonCrystals || c.inv.PrizeCount() == 10)
	case "magicExtensions":
		bars := (b2i(c.inv.HalfMagic()) + 1) * (b2i(c.inv.QuarterMagic()) + 1)
		return bars * (c.inv.BottleCount() + 1)
	}
	_ = arg
	return 0
}
