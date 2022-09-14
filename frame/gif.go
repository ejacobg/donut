package frame

import (
	"image"
	"image/color"
)

// GIF implements the Frame interface.
type GIF struct {
	Image image.Paletted
	SF    SetFunc[color.Color]
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
