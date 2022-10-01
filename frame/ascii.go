package frame

import (
	"fmt"
	"io"
)

// ASCII implements the Frame interface.
type ASCII struct {
	// A set of luminance values, represented as characters.
	frame [][]rune
	SF    SetFunc[rune]
}

func makeFrame(width, height int) [][]rune {
	frame := make([][]rune, height)
	for i := range frame {
		frame[i] = make([]rune, width)
	}
	return frame
}

func NewASCII(width, height int, sf SetFunc[rune]) *ASCII {
	frame := makeFrame(width, height)
	return &ASCII{frame, sf}
}

func (a *ASCII) Set(x, y int, L float64) bool {
	if a.SF == nil {
		return false
	}
	if r, ok := a.SF(L); ok {
		a.frame[x][y] = r
		return true
	}
	return false
}

func (a *ASCII) SetAll(r rune) {
	for i := range a.frame {
		for j := range a.frame[i] {
			a.frame[i][j] = r
		}
	}
}

func (a *ASCII) Reset() {
	a.SetAll(' ')
}

func (a *ASCII) Resize(width int, height int) {
	a.frame = makeFrame(width, height)
}

func (a *ASCII) Print(w io.Writer) {
	for i := range a.frame {
		_, _ = fmt.Fprintln(w, string(a.frame[i]))
	}
}

// ASCIISF is the original SetFunc as described in the original implementation.
// Assumes L ranges from -sqrt(2) to +sqrt(2).
func ASCIISF(L float64) (rune, bool) {
	// L ranges from -sqrt(2) to +sqrt(2).  If it's < 0, the surface
	// is pointing away from us, so we won't bother trying to plot it.
	if L > 0 {
		luminanceIndex := int(L * 8)
		// luminanceIndex is now in the range 0..11 (8*sqrt(2) = 11.3)
		// now we look up the character corresponding to the
		// luminance and plot it in our output:
		return []rune(".,-~:;=!*#$@")[luminanceIndex], true
	}
	return 0, false
}
