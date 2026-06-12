package maps

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"sync"
)

//go:embed assets/lightworld.png
var lightWorldPNG []byte

//go:embed assets/darkworld.png
var darkWorldPNG []byte

// DisplaySize is the initial rendered width and height of each world map,
// in pixels; the maps rescale with the window afterwards.
const DisplaySize = 320

type world struct {
	name string
	data []byte

	once sync.Once
	src  *image.NRGBA
	err  error
}

var (
	lightWorld = &world{name: "lightworld", data: lightWorldPNG}
	darkWorld  = &world{name: "darkworld", data: darkWorldPNG}
)

// LightWorld returns the light world map scaled to a size x size square.
func LightWorld(size int) (image.Image, error) { return lightWorld.scaled(size) }

// DarkWorld returns the dark world map scaled to a size x size square.
func DarkWorld(size int) (image.Image, error) { return darkWorld.scaled(size) }

func (w *world) scaled(size int) (image.Image, error) {
	w.once.Do(func() {
		src, err := png.Decode(bytes.NewReader(w.data))
		if err != nil {
			w.err = fmt.Errorf("maps: decoding %s: %w", w.name, err)
			return
		}
		b := src.Bounds()
		w.src = image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
		draw.Draw(w.src, w.src.Bounds(), src, b.Min, draw.Src)
	})
	if w.err != nil {
		return nil, w.err
	}
	return scaleTo(w.src, size), nil
}

// scaleTo box-filters src down to a size x size square. Averaging the source
// region per output pixel avoids the aliasing nearest-neighbor would produce
// at large reduction ratios.
func scaleTo(src *image.NRGBA, size int) *image.NRGBA {
	w, h := src.Bounds().Dx(), src.Bounds().Dy()
	out := image.NewNRGBA(image.Rect(0, 0, size, size))
	for y := 0; y < size; y++ {
		sy0 := y * h / size
		sy1 := max((y+1)*h/size, sy0+1)
		for x := 0; x < size; x++ {
			sx0 := x * w / size
			sx1 := max((x+1)*w/size, sx0+1)
			var r, g, bl, a, n int
			for sy := sy0; sy < sy1; sy++ {
				row := src.Pix[sy*src.Stride+sx0*4 : sy*src.Stride+sx1*4]
				for i := 0; i < len(row); i += 4 {
					r += int(row[i])
					g += int(row[i+1])
					bl += int(row[i+2])
					a += int(row[i+3])
					n++
				}
			}
			o := out.PixOffset(x, y)
			out.Pix[o] = uint8(r / n)
			out.Pix[o+1] = uint8(g / n)
			out.Pix[o+2] = uint8(bl / n)
			out.Pix[o+3] = uint8(a / n)
		}
	}
	return out
}
