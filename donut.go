package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/quincinia/donut/frame"
	"golang.org/x/term"
)

const frameDelay = 50 // milliseconds

func main() {
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalln(err)
	}

	var min int
	if width < height {
		min = width
	} else {
		// If we're constrained by height, we subtract 1 because fmt.Println adds an extra line.
		// If we don't do this, we get artifacts when multiple renders occur.
		min = height - 1
	}

	ascii := frame.NewASCII(min, min, nil)
	// The orginal SetFunc as described in the original implementation.
	// Assumes L ranges from -sqrt(2) to +sqrt(2).
	ascii.SF = frame.SetFunc[rune](func(L float64) (rune, bool) {
		// L ranges from -sqrt(2) to +sqrt(2).  If it's < 0, the surface
		// is pointing away from us, so we won't bother trying to plot it.
		if L > 0 {
			luminanceIndex := int(L * 8)
			// luminanceIndex is now in the range 0..11 (8*sqrt(2) = 11.3)
			// now we lookup the character corresponding to the
			// luminance and plot it in our output:
			return []rune(".,-~:;=!*#$@")[luminanceIndex], true
		}
		return 0, false
	})
	torus := frame.DefaultTorus()
	scene := frame.Scene{
		Width:  min,
		Height: min,
		K2:     5.0,
		LX:     0.0,
		LY:     1.0,
		LZ:     -1.0,
	}
	scene.CalculateK1(torus.R1, torus.R2)

	// The number of frames to render in a full cycle.
	steps := 200

	// How much the angles should change between each step.
	startA, startB := 15.0, 25.0
	stepA, stepB := 0.07, 0.03

	for {
		var s int
		for s, torus.A, torus.B = 0, startA, startB; s < steps; s, torus.A, torus.B = s+1, torus.A+stepA, torus.B+stepB {
			ascii.SetAll(' ') // Should the Render function clear the frame instead of doing it here?
			frame.Render(ascii, torus, scene)
			ascii.Print(os.Stdout)
			fmt.Printf("\033[%dF", min)
			time.Sleep(frameDelay * time.Millisecond)
		}
	}
}
