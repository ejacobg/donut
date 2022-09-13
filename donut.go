package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"golang.org/x/term"
)

const (
	thetaSpacing  = 0.07
	phiSpacing    = 0.02
	frameDelay    = 50 // milliseconds
)

const (
	R1 = 1.0
	R2 = 2.0
	K2 = 5.0
)

// Experimenting with generics.
func newMatrix[T any](dim int) *[][]T {
	matrix := make([][]T, dim)
	for i := range matrix {
		matrix[i] = make([]T, dim)
	}
	return &matrix
}

var (
	newFrame   = newMatrix[rune]
	newZBuffer = newMatrix[float64]
)

// renderFrame writes a frame of the animation to stdout.
// A and B are the rotation angles (in radians) around the X and Z axes, respectively.
func renderFrame(A, B float64, dim int) {
	cosA, sinA := math.Cos(A), math.Sin(A)
	cosB, sinB := math.Cos(B), math.Sin(B)

	frame, zbuffer := newFrame(dim), newZBuffer(dim)
	for i := range *frame {
		for j := range (*frame)[i] {
			(*frame)[i][j] = ' '
		}
	}

	K1 := float64(dim) * K2 * 3 / (8 * (R1 + R2))

	// theta goes around the cross-sectional circle of a torus
	for theta := 0.0; theta < 2.0*math.Pi; theta += thetaSpacing {
		cosTheta, sinTheta := math.Cos(theta), math.Sin(theta)

		// phi goes around the center of revolution of a torus
		for phi := 0.0; phi < 2.0*math.Pi; phi += phiSpacing {
			cosPhi, sinPhi := math.Cos(phi), math.Sin(phi)

			// the x,y coordinate of the circle, before revolving (factored out of the above equations)
			circleX := R2 + R1*cosTheta
			circleY := R1 * sinTheta

			// final 3D (x,y,z) coordinate after rotations, directly from our math above
			x := circleX*(cosB*cosPhi+sinA*sinB*sinPhi) - circleY*cosA*sinB
			y := circleX*(sinB*cosPhi-sinA*cosB*sinPhi) + circleY*cosA*cosB
			z := K2 + cosA*circleX*sinPhi + circleY*sinA
			ooz := 1 / z // "one over z"

			// x and y projection.  note that y is negated here, because y
			// goes up in 3D space but down on 2D displays.
			xp := int(float64(dim)/2 + K1*ooz*x)
			yp := int(float64(dim)/2 - K1*ooz*y)

			// calculate luminance.  ugly, but correct.
			L := cosPhi*cosTheta*sinB - cosA*cosTheta*sinPhi - sinA*sinTheta + cosB*(cosA*sinTheta-cosTheta*sinA*sinPhi)
			// L ranges from -sqrt(2) to +sqrt(2).  If it's < 0, the surface
			// is pointing away from us, so we won't bother trying to plot it.
			if L > 0 {
				// test against the z-buffer.  larger 1/z means the pixel is
				// closer to the viewer than what's already plotted.
				if ooz > (*zbuffer)[xp][yp] {
					(*zbuffer)[xp][yp] = ooz
					luminanceIndex := int(L * 8)
					// luminanceIndex is now in the range 0..11 (8*sqrt(2) = 11.3)
					// now we lookup the character corresponding to the
					// luminance and plot it in our output:
					(*frame)[xp][yp] = []rune(".,-~:;=!*#$@")[luminanceIndex]
				}
			}
		}
	}

	for i := range *frame {
		fmt.Println(string((*frame)[i]))
	}
	fmt.Printf("\033[%dF", dim)
}

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
		min = height-1
	}

	// The number of frames to render in a full cycle.
	steps := 100

	// How much the angles should change between each step.
	startA, startB := 15.0, 25.0
	stepA, stepB := 0.07, 0.03

	for {
		for s, A, B := 0, startA, startB; s < steps; s, A, B = s+1, A+stepA, B+stepB {
			renderFrame(A, B, min)
			time.Sleep(frameDelay * time.Millisecond)
		}
	}
}
