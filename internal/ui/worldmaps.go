package ui

import (
	"image"
	"image/color"
	"image/draw"
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

	// Markers overlaid on each map, in normalized coordinates. Redrawn into
	// the scaled images on every render.
	lightMarks, darkMarks []MapMarker

	// Debounce state for interactive resizes. Configure callbacks and the
	// TclAfter they schedule both run on the Tk event loop, so these need no
	// locking.
	pendingAfter   string
	allocW, allocH int
}

// MapMarker is a normalized [0,1] position on a map together with whether the
// location there is currently reachable in logic (which colors it).
type MapMarker struct {
	U, V       float64
	Accessible bool
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
	w.size = size
	return w.render()
}

// SetMarkers replaces the overlaid location markers and redraws. Must be
// called on the Tk event loop.
func (w *WorldMaps) SetMarkers(light, dark []MapMarker) error {
	w.lightMarks = light
	w.darkMarks = dark
	return w.render()
}

// render rescales both maps to the current size, paints the markers onto them,
// and swaps in fresh photos (deleting the old ones so Tk frees them).
func (w *WorldMaps) render() error {
	light, err := maps.LightWorld(w.size)
	if err != nil {
		return err
	}
	dark, err := maps.DarkWorld(w.size)
	if err != nil {
		return err
	}
	drawMarkers(light, w.lightMarks)
	drawMarkers(dark, w.darkMarks)

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
	return nil
}

var (
	markerInLogic    = color.NRGBA{R: 30, G: 215, B: 50, A: 255}  // vivid green
	markerOutOfLogic = color.NRGBA{R: 230, G: 40, B: 40, A: 255}  // vivid red
	markerBorder     = color.NRGBA{R: 0, G: 0, B: 0, A: 255}      // opaque black
)

// drawMarkers paints a square per marker onto img, which maps.LightWorld/
// DarkWorld produce as a mutable *image.NRGBA. Coordinates are normalized, so
// markers land correctly at any map size.
func drawMarkers(img image.Image, marks []MapMarker) {
	dst, ok := img.(draw.Image)
	if !ok {
		return
	}
	b := img.Bounds()
	size := b.Dx()
	// Half-extent and border scale with the map but stay legible when small;
	// the thick black frame is what makes the colors readable over the map.
	half := min(max(size/80, 4), 12)
	border := min(max(size/300, 2), 4)
	for _, m := range marks {
		cx := b.Min.X + int(m.U*float64(size))
		cy := b.Min.Y + int(m.V*float64(size))
		fill := markerOutOfLogic
		if m.Accessible {
			fill = markerInLogic
		}
		drawSquare(dst, cx, cy, half, border, fill)
	}
}

// drawSquare fills a square of half-extent `half` centered at (cx,cy), framed
// by a black border `border` pixels thick. dst.Set is bounds-checked, so
// squares near an edge are simply clipped.
func drawSquare(dst draw.Image, cx, cy, half, border int, fill color.NRGBA) {
	for dy := -half; dy <= half; dy++ {
		for dx := -half; dx <= half; dx++ {
			edge := min(min(dx+half, half-dx), min(dy+half, half-dy))
			c := fill
			if edge < border {
				c = markerBorder
			}
			dst.Set(cx+dx, cy+dy, c)
		}
	}
}

// atoiOr parses s as an int, returning def if it isn't a valid number. Tk
// reports event geometry fields as strings.
func atoiOr(s string, def int) int {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return def
}
