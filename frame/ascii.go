package frame

import (
	"fmt"
	"io"
)

// ASCII implements the Frame interface.
type ASCII struct {
	// A set of luminance values, represented as characters.
	frame [][]rune
	SF    SetFunc
}

func makeFrame(width, height int) [][]rune {
	frame := make([][]rune, height)
	for i := range frame {
		frame[i] = make([]rune, width)
	}
	return frame
}

func NewASCII(width, height int, sf SetFunc) *ASCII {
	frame := makeFrame(width, height)
	return &ASCII{frame, sf}
}

func (a *ASCII) Set(x, y int, L float64) {
	if a.SF == nil {
		return
	}
	a.SF(x, y, L)
}

func (a *ASCII) SetAll(r rune) {
	for i := range a.frame {
		for j := range a.frame[i] {
			a.frame[i][j] = r
		}
	}
}

func (a *ASCII) Resize(width int, height int) {
	a.frame = makeFrame(width, height)
}

func (a *ASCII) Print(w io.Writer) {
	for i := range a.frame {
		fmt.Fprintln(w, string(a.frame[i]))
	}
}
