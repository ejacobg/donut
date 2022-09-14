package frame

// ASCII implements the Frame interface.
type ASCII struct {
	frame [][]rune
	SF    SetFunc
}

func NewASCII(width, height int, sf SetFunc) *ASCII {
	frame := make([][]rune, height)
	for i := range frame {
		frame[i] = make([]rune, width)
	}
	return &ASCII{frame, sf}
}

func (a *ASCII) Set(x, y int, L float64) {
	if a.SF == nil {
		return
	}
	a.SF(x, y, L)
}
