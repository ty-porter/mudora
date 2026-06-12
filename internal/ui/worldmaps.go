package ui

import (
	"github.com/ty-porter/mudora/internal/ui/maps"
	. "modernc.org/tk9.0"
)

const (
	// minMapSize keeps the maps legible when the window gets small;
	// maxMapSize avoids upscaling past the source's useful detail.
	minMapSize = 200
	maxMapSize = 1000
	// resizeStep ignores sub-step size changes; rate limiting during
	// interactive resizes is handled by the caller's debounce, so this
	// only suppresses layout jitter.
	resizeStep = 4
)

// WorldMaps displays the light and dark world maps side by side, rescaling
// to fit the window.
type WorldMaps struct {
	frame *TFrameWidget
	light *LabelWidget
	dark  *LabelWidget
	size  int
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

	w.light = w.frame.Label(Image(NewPhoto(Data(light))), Padx(0), Pady(0), Borderwidth(0))
	w.dark = w.frame.Label(Image(NewPhoto(Data(dark))), Padx(0), Pady(0), Borderwidth(0))
	Grid(w.light, Row(0), Column(0), Padx("4"))
	Grid(w.dark, Row(0), Column(1), Padx("4"))

	return w, nil
}

// Frame returns the container widget for placement in a parent layout.
func (w *WorldMaps) Frame() *TFrameWidget {
	return w.frame
}

// Resize rescales both maps to size x size each. Changes smaller than
// resizeStep are ignored.
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
	w.size = size
	w.light.Configure(Image(NewPhoto(Data(light))))
	w.dark.Configure(Image(NewPhoto(Data(dark))))
	return nil
}
