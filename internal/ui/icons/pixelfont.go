package icons

import (
	"image"
	"image/color"
	"image/draw"
)

// pixelGlyphs is a tiny bitmap font for cell captions, drawn at native pixel
// scale. '#' pixels are lit. Glyphs may vary in width; all are glyphHeight
// tall.
var pixelGlyphs = map[rune][]string{
	'D': {"### ", "#  #", "#  #", "#  #", "#  #", "### "},
	'E': {"####", "#   ", "### ", "#   ", "#   ", "####"},
	'H': {"#  #", "#  #", "####", "#  #", "#  #", "#  #"},
	'I': {"###", " # ", " # ", " # ", " # ", "###"},
	'M': {"#   #", "## ##", "# # #", "#   #", "#   #", "#   #"},
	'P': {"### ", "#  #", "#  #", "### ", "#   ", "#   "},
	'R': {"### ", "#  #", "#  #", "### ", "# # ", "#  #"},
	'S': {" ###", "#   ", " ## ", "   #", "   #", "### "},
	'T': {"###", " # ", " # ", " # ", " # ", " # "},
	'W': {"#   #", "#   #", "#   #", "# # #", "# # #", " # # "},
	'o': {"    ", "    ", " ## ", "#  #", "#  #", " ## "},
}

const glyphHeight = 6

// captionColor matches the light marker color used on sprites.
var captionColor = color.NRGBA{R: 230, G: 230, B: 230, A: 255}

// textWidth returns the unscaled pixel width of s.
func textWidth(s string) int {
	w := 0
	for _, r := range s {
		if g, ok := pixelGlyphs[r]; ok {
			w += len(g[0]) + 1
		}
	}
	return max(w-1, 0)
}

// TextWidth returns the display pixel width of s. Captions draw at native
// scale (not upscaled with the sprites), so this equals the glyph width.
func TextWidth(s string) int { return textWidth(s) }

// drawText draws s into dst with its top-left at (x, y), unscaled.
func drawText(dst *image.NRGBA, x, y int, s string) {
	for _, r := range s {
		g, ok := pixelGlyphs[r]
		if !ok {
			continue
		}
		for gy, row := range g {
			for gx, c := range row {
				if c == '#' {
					dst.SetNRGBA(x+gx, y+gy, captionColor)
				}
			}
		}
		x += len(g[0]) + 1
	}
}

// Captioned returns sprite with caption drawn to its left. zoneWidth is the
// display width reserved for the caption — pass the TextWidth of the longest
// caption in a column so its sprites align. The caption is right-aligned and
// the sprite left-aligned, keeping the pair tight; leftPad separates the
// cell from any neighboring grid column.
func Captioned(caption string, zoneWidth int, sprite image.Image) image.Image {
	const (
		leftPad = 8 // display px of breathing room before the caption
		gap     = 2 // display px between caption and sprite
	)
	cell := 16 * scaleNum / scaleDen
	sb := sprite.Bounds()
	h := max(cell, sb.Dy())
	out := image.NewNRGBA(image.Rect(0, 0, leftPad+zoneWidth+gap+cell, h))

	if tw := textWidth(caption); tw > 0 {
		text := image.NewNRGBA(image.Rect(0, 0, tw, glyphHeight))
		drawText(text, 0, 0, caption)
		x0 := leftPad + zoneWidth - tw
		y0 := (h - glyphHeight) / 2
		draw.Draw(out, image.Rect(x0, y0, x0+tw, y0+glyphHeight), text, image.Point{}, draw.Over)
	}

	x0 := leftPad + zoneWidth + gap
	y0 := (h - sb.Dy()) / 2
	draw.Draw(out, image.Rect(x0, y0, x0+sb.Dx(), y0+sb.Dy()), sprite, sb.Min, draw.Over)
	return out
}
