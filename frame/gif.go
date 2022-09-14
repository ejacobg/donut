package frame

import "image"

// GIF implements the Frame interface.
type GIF struct {
	Image image.Paletted
	SF    SetFunc
}

func (a *GIF) Set(x, y int, L float64) {
	if a.SF == nil {
		return
	}
	a.SF(x, y, L)
}
