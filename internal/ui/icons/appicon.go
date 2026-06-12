package icons

import (
	"image"
	"image/draw"
)

var appIconBox = Box{X: 126, Y: 121, W: 13, H: 15}

func AppIcon() (image.Image, error) {
	s, err := loadSheet()
	if err != nil {
		return nil, err
	}
	b := appIconBox
	// Upscale twice (4x total) so the window manager downscales a large
	// source instead of upscaling a tiny one.
	sprite := upscale(upscale(s.SubImage(image.Rect(b.X, b.Y, b.X+b.W, b.Y+b.H))))
	sb := sprite.Bounds()
	side := max(sb.Dx(), sb.Dy())
	out := image.NewNRGBA(image.Rect(0, 0, side, side))
	x0 := (side - sb.Dx()) / 2
	y0 := (side - sb.Dy()) / 2
	draw.Draw(out, image.Rect(x0, y0, x0+sb.Dx(), y0+sb.Dy()), sprite, sb.Min, draw.Over)
	return out, nil
}
