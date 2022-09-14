package frame

import (
	"image"
	"image/color"
	"image/gif"
	"io"
)

// GIF implements the Frame interface.
type GIF struct {
	Image   *image.Paletted
	Options gif.Options
	SF      SetFunc[color.Color]
}

func (g *GIF) Set(x, y int, L float64) bool {
	if g.SF == nil {
		return false
	}
	if c, ok := g.SF(L); ok {
		g.Image.Set(x, y, c)
		return true
	}
	return false
}

// The reset color is the first item in the Palette.
func (g *GIF) Reset() {
	for i := range g.Image.Pix {
		g.Image.Pix[i] = 0
	}
}

func (g *GIF) Print(w io.Writer) {
	gif.Encode(w, g.Image, &g.Options)
}

// 15 Gray16 colors plus a transparent color.
func Gray16Palette() color.Palette {
	palette := color.Palette{color.Alpha{0}}
	for i := 0; i < 16; i++ {
		palette = append(palette, color.Gray16{1 << i})
	}
	return palette
}

// Assumes L ranges from -sqrt(2) to +sqrt(2).
func Gray16SF(L float64) (color.Color, bool) {
	// index ranges from [1, 15]
	index := int(L * 5) + 8
	return color.Gray16{1<<index}, true
}
