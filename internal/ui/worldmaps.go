package ui

import (
	"strconv"
	"time"

	"github.com/ty-porter/mudora/internal/ui/maps"
	. "modernc.org/tk9.0"
)

const (
	// minMapSize keeps the maps legible when the window gets small;
	// maxMapSize avoids upscaling past the source's useful detail.
	minMapSize = 200
	maxMapSize = 1000
	// resizeStep ignores sub-step size changes so interactive resizes don't
	// rescale (an O(source-pixels) box filter) on every pixel of drag, and so
	// the size can't oscillate by a pixel between layout passes.
	resizeStep = 4
	// resizeDebounce delays the rescale until a drag settles: each Configure
	// re-arms the timer, so the costly rescale runs once the user stops.
	resizeDebounce = 100 * time.Millisecond
	// mapPadX is the total horizontal padding around the two maps (Padx("4")
	// on each side of each map), subtracted from the width budget.
	mapPadX = 16
)

// WorldMaps displays the light and dark world maps side by side, rescaling to
// fit the space the layout gives its frame.
type WorldMaps struct {
	frame *TFrameWidget
	light *LabelWidget
	dark  *LabelWidget
	// Photos are kept so each can be deleted when replaced; otherwise Tk would
	// orphan a photo image on every rescale.
	lightPhoto *Img
	darkPhoto  *Img
	size       int

	// Debounce state for interactive resizes. Configure callbacks and the
	// TclAfter they schedule both run on the Tk event loop, so these need no
	// locking.
	pendingAfter   string
	allocW, allocH int
}

func NewWorldMaps() (*WorldMaps, error) {
	w := &WorldMaps{frame: TFrame(), size: maps.DisplaySize}

	light, err := maps.LightWorld(w.size)
	if err != nil {
		return nil, err
	}
	dark, err := maps.DarkWorld(w.size)
	if err != nil {
		return nil, err
	}

	w.lightPhoto = NewPhoto(Data(light))
	w.darkPhoto = NewPhoto(Data(dark))
	w.light = w.frame.Label(Image(w.lightPhoto), Padx(0), Pady(0), Borderwidth(0))
	w.dark = w.frame.Label(Image(w.darkPhoto), Padx(0), Pady(0), Borderwidth(0))
	// Sticky E/W with equal column weights keeps the two maps centered as a
	// pair when the frame is stretched wider than they are.
	Grid(w.light, Row(0), Column(0), Padx("4"), Sticky(E))
	Grid(w.dark, Row(0), Column(1), Padx("4"), Sticky(W))
	GridColumnConfigure(w.frame, 0, Weight(1))
	GridColumnConfigure(w.frame, 1, Weight(1))
	GridRowConfigure(w.frame, 0, Weight(1))

	// Rescale to fit whenever the layout resizes the frame. Bound to the frame
	// itself, so only the frame's own Configure fires this (not its children's
	// or the toplevel's). Interactive drags fire rapidly, so the rescale is
	// debounced.
	Bind(w.frame, "<Configure>", Command(w.onConfigure))

	return w, nil
}

// Frame returns the container widget for placement in a parent layout.
func (w *WorldMaps) Frame() *TFrameWidget {
	return w.frame
}

// onConfigure records the frame's new allocated size and (re)arms the
// debounced rescale. Runs on the Tk event loop.
func (w *WorldMaps) onConfigure(e *Event) {
	w.allocW = atoiOr(e.Width, w.allocW)
	w.allocH = atoiOr(e.Height, w.allocH)
	if w.pendingAfter != "" {
		TclAfterCancel(w.pendingAfter)
	}
	w.pendingAfter = TclAfter(resizeDebounce, w.applyResize)
}

// applyResize fits the maps to the frame's last allocated size. The two maps
// sit side by side, so each gets half the width (less the padding between and
// around them); each is square, so the usable size is bounded by height too.
func (w *WorldMaps) applyResize() {
	w.pendingAfter = ""
	widthBudget := (w.allocW - mapPadX) / 2
	heightBudget := w.allocH
	// Resize clamps to [minMapSize, maxMapSize] and ignores sub-step changes.
	_ = w.Resize(min(widthBudget, heightBudget))
}

// Resize rescales both maps to size x size each. The size is clamped to the
// legible range, and changes smaller than resizeStep are ignored.
func (w *WorldMaps) Resize(size int) error {
	size = min(max(size, minMapSize), maxMapSize)
	d := size - w.size
	if -resizeStep < d && d < resizeStep {
		return nil
	}

	light, err := maps.LightWorld(size)
	if err != nil {
		return err
	}
	dark, err := maps.DarkWorld(size)
	if err != nil {
		return err
	}

	lightPhoto := NewPhoto(Data(light))
	darkPhoto := NewPhoto(Data(dark))
	w.light.Configure(Image(lightPhoto))
	w.dark.Configure(Image(darkPhoto))
	// Swap in the new photos before deleting the old ones, which now have no
	// displayed instances and so are freed immediately.
	w.lightPhoto.Delete()
	w.darkPhoto.Delete()
	w.lightPhoto = lightPhoto
	w.darkPhoto = darkPhoto
	w.size = size
	return nil
}

// atoiOr parses s as an int, returning def if it isn't a valid number. Tk
// reports event geometry fields as strings.
func atoiOr(s string, def int) int {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return def
}
