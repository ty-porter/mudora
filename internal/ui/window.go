package ui

import (
	"context"
	"sync"
	"time"

	"github.com/ty-porter/mudora/internal/game"
	"github.com/ty-porter/mudora/internal/ui/icons"
	. "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

// uiTick is how often the Tk goroutine drains state pushed from the tracker.
// Updates are coalesced, so this only bounds redraw latency, not throughput.
const uiTick = 100 * time.Millisecond

// statusWidth caps the status bar at a fixed character budget so long
// messages (e.g. verbose device URIs) can't push the window wider than its
// natural map-driven width. It sits comfortably below that width, so the
// status line never drives the layout.
const statusWidth = 48

// TrackerWindow owns the Tk widgets and the bridge between the tracker's poll
// goroutine and the Tk event loop. ApplyState/SetStatus are safe to call from
// any goroutine; they stash the latest value, which a TclAfter tick on the Tk
// goroutine applies. No Tk calls are ever made off the Tk goroutine.
type TrackerWindow struct {
	mu        sync.Mutex
	pending   *game.State
	status    string
	statusNew bool

	ctx         context.Context
	items       *ItemGrid
	statusLabel *TLabelWidget
}

func New() *TrackerWindow {
	return &TrackerWindow{}
}

// ApplyState records the latest game state for the UI to render on its next
// tick. Safe to call from any goroutine.
func (w *TrackerWindow) ApplyState(s game.State) {
	w.mu.Lock()
	w.pending = &s
	w.mu.Unlock()
}

// SetStatus records a status message for the status bar. Safe to call from any
// goroutine.
func (w *TrackerWindow) SetStatus(msg string) {
	w.mu.Lock()
	w.status = msg
	w.statusNew = true
	w.mu.Unlock()
}

// drain runs on the Tk goroutine: it applies any pending state/status, redraws,
// then re-arms itself. This is the only place tracker data reaches the widgets.
func (w *TrackerWindow) drain() {
	// The context is cancelled when the process is interrupted (e.g. Ctrl+C in
	// the terminal). Tk's event loop won't notice on its own, so tear down the
	// window here, which unblocks App.Wait and lets the process exit.
	if w.ctx.Err() != nil {
		Destroy(App)
		return
	}

	w.mu.Lock()
	pending := w.pending
	w.pending = nil
	status, statusNew := w.status, w.statusNew
	w.statusNew = false
	w.mu.Unlock()

	if pending != nil {
		applyState(*pending)
		w.items.RefreshAll()
	}
	if statusNew {
		w.statusLabel.Configure(Txt(truncateMiddle(status, statusWidth)))
	}

	TclAfter(uiTick, w.drain)
}

func (w *TrackerWindow) Run(ctx context.Context) error {
	w.ctx = ctx
	if err := ActivateTheme("azure dark"); err != nil {
		return err
	}
	if bg := StyleLookup("TFrame", Background); bg != "" {
		App.Configure(Background(bg))
	}
	App.WmTitle("Mudora")
	appIcon, err := icons.AppIcon()
	if err != nil {
		return err
	}
	App.IconPhoto(NewPhoto(Data(appIcon)))

	worldMaps, err := NewWorldMaps()
	if err != nil {
		return err
	}
	notes := NewNotes()
	items, err := NewItemGrid()
	if err != nil {
		return err
	}
	rewards, err := NewDungeonRewardGrid()
	if err != nil {
		return err
	}
	w.items = items
	// Width pins the label's requested size so overflow is clipped, not grown
	// into; truncateMiddle keeps the displayed text within that same budget.
	w.statusLabel = TLabel(Txt("starting…"), Width(statusWidth))
	// World maps in row 0 spanning all columns, section headers in row 1,
	// content in row 2, separators spanning headers and content.
	// Sticky so the frame fills its cell; its Configure then tracks the window
	// size and drives the map rescale.
	Grid(worldMaps.Frame(), Row(0), Column(0), Columnspan(5), Sticky(NEWS), Pady("4"))

	Grid(TLabel(Txt("Notes")), Row(1), Column(0), Sticky(W), Padx("4"), Pady("4"))
	Grid(TLabel(Txt("Inventory")), Row(1), Column(2), Sticky(W), Padx("4"), Pady("4"))
	Grid(TLabel(Txt("Dungeon Rewards")), Row(1), Column(4), Sticky(W), Padx("4"), Pady("4"))

	Grid(notes.Widget(), Row(2), Column(0), Sticky(NEWS), Padx("4"), Pady("4"))
	Grid(TSeparator(Orient("vertical")), Row(1), Column(1), Rowspan(2), Sticky(NS), Pady("4"))
	Grid(items.Frame(), Row(2), Column(2), Sticky(N), Padx("4"), Pady("4"))
	Grid(TSeparator(Orient("vertical")), Row(1), Column(3), Rowspan(2), Sticky(NS), Pady("4"))
	Grid(rewards.Frame(), Row(2), Column(4), Sticky(N), Padx("4"), Pady("4"))
	// Status bar spans the full width below the content.
	Grid(TSeparator(Orient("horizontal")), Row(3), Column(0), Columnspan(5), Sticky(WE), Pady("4"))
	Grid(w.statusLabel, Row(4), Column(0), Columnspan(5), Sticky(W), Padx("4"), Pady("2"))
	// Extra window space goes to the notes column; the icon grids keep
	// their natural size.
	GridColumnConfigure(App, 0, Weight(1))
	GridRowConfigure(App, 0, Weight(1))

	// Start draining tracker updates on the Tk event loop, then enter it.
	TclAfter(uiTick, w.drain)
	App.Wait()

	return nil
}

// truncateMiddle shortens s to at most max runes, replacing the excised
// middle with an ellipsis so both ends stay visible (status lines carry the
// device on the left and the tracking state on the right).
func truncateMiddle(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	if max <= 1 {
		return string(r[:max])
	}
	keep := max - 1 // one rune spent on the ellipsis
	head := (keep + 1) / 2
	tail := keep - head
	return string(r[:head]) + "…" + string(r[len(r)-tail:])
}
