package icons

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"sync"
)

//go:embed assets/sprites.png
var sheetPNG []byte

// The sheet has no alpha channel; this background color is keyed to
// transparent when the sheet is first loaded.
const keyR, keyG, keyB = 128, 0, 128

// Box locates a sprite on the sheet: the offset of its top-left corner and
// its size, in pixels.
type Box struct {
	X, Y, W, H int
}

// GridItem is one cell of a tracker grid. Progressive items carry one
// bounding box per upgrade tier. State 0 means the item is not yet obtained
// and draws the first tier's sprite grayed out; state N draws tier N (the
// sprite at boundingBoxes[N-1]) in full color.
type GridItem struct {
	boundingBoxes []Box
	State         int
}

// Icon returns the sprite for the item's current state, upscaled for
// display. Disabled items (state 0) are grayed out.
func (it *GridItem) Icon() (image.Image, error) {
	if it.State < 0 || it.State > len(it.boundingBoxes) {
		return nil, fmt.Errorf("icons: state %d out of range (%d tiers)", it.State, len(it.boundingBoxes))
	}
	s, err := loadSheet()
	if err != nil {
		return nil, err
	}
	b := it.boundingBoxes[max(it.State-1, 0)]
	crop := image.Image(s.SubImage(image.Rect(b.X, b.Y, b.X+b.W, b.Y+b.H)))
	if it.State == 0 {
		crop = grayscale(crop)
	}
	return upscale(crop), nil
}

var (
	sheetOnce sync.Once
	sheet     *image.NRGBA
	sheetErr  error
)

func loadSheet() (*image.NRGBA, error) {
	sheetOnce.Do(func() {
		src, err := png.Decode(bytes.NewReader(sheetPNG))
		if err != nil {
			sheetErr = fmt.Errorf("icons: decoding spritesheet: %w", err)
			return
		}
		img := image.NewNRGBA(src.Bounds())
		draw.Draw(img, img.Bounds(), src, src.Bounds().Min, draw.Src)
		for i := 0; i < len(img.Pix); i += 4 {
			if img.Pix[i] == keyR && img.Pix[i+1] == keyG && img.Pix[i+2] == keyB {
				img.Pix[i+3] = 0
			}
		}
		sheet = img
	})
	return sheet, sheetErr
}

// Sheet sprites are SNES-native (~16px)
const scaleNum, scaleDen = 2, 1

func upscale(src image.Image) image.Image {
	b := src.Bounds()
	w, h := b.Dx()*scaleNum/scaleDen, b.Dy()*scaleNum/scaleDen
	out := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		sy := b.Min.Y + y*scaleDen/scaleNum
		for x := 0; x < w; x++ {
			sx := b.Min.X + x*scaleDen/scaleNum
			out.Set(x, y, src.At(sx, sy))
		}
	}
	return out
}

// grayscale returns a dimmed grayscale copy of src, preserving alpha.
func grayscale(src image.Image) image.Image {
	b := src.Bounds()
	out := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(out, out.Bounds(), src, b.Min, draw.Src)
	for i := 0; i < len(out.Pix); i += 4 {
		r, g, bl := int(out.Pix[i]), int(out.Pix[i+1]), int(out.Pix[i+2])
		lum := (299*r + 587*g + 114*bl) / 1000
		lum = lum * 3 / 4
		out.Pix[i], out.Pix[i+1], out.Pix[i+2] = uint8(lum), uint8(lum), uint8(lum)
	}
	return out
}
