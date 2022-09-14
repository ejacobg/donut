package frame

import "math"

type Frame interface {
	// Apply a luminance value to a point in the frame.
	Set(x, y int, L float64)
}

// Properties of the torus to be rendered.
type Torus struct {
	// Controls the resolution of the cross-sectional circle.
	ThetaSpacing float64

	// Controls the resolution of the revolution of the circle.
	PhiSpacing float64

	// Radius of the cross-sectional circle.
	R1 float64

	// Distance from axis of revolution to the center of the cross-sectional circle.
	R2 float64

	// Rotation of the torus along the X and Z axes.
	A, B float64
}

func DefaultTorus() Torus {
	return Torus{
		ThetaSpacing: 0.07,
		PhiSpacing:   0.02,
		R1:           1,
		R2:           2,
	}
}

// Locations and properties of the scene.
type Scene struct {
	// Dimensions of the viewport.
	Width, Height int

	// Distance from camera to viewport.
	K1 float64

	// Location of the torus with respect to the center of the cross-sectional circle.
	TX, TY, TZ float64

	// Distance from camera to torus. Essentially equal to TZ since camera is at origin.
	K2 float64

	// Location of the light source.
	LX, LY, LZ float64
}

// Renders the scene onto the frame. Assumes that the frame is sized correctly.
func Render(frame Frame, t Torus, s Scene) {
	cosA, sinA := math.Cos(A), math.Sin(A)
	cosB, sinB := math.Cos(B), math.Sin(B)

	// The Palleted pixel colors are by default set to the 0 index of the Palette slice, so no initialization is needed here.
	zbuffer := rectZBuffer(width, height)

	K1 := float64(width) * K2 * 3 / (8 * (R1 + R2))

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
			xp := int(float64(width)/2 + K1*ooz*x)
			yp := int(float64(height)/2 - K1*ooz*y)

			// calculate luminance.  ugly, but correct.
			L := cosPhi*cosTheta*sinB - cosA*cosTheta*sinPhi - sinA*sinTheta + cosB*(cosA*sinTheta-cosTheta*sinA*sinPhi)
			// L ranges from -sqrt(2) to +sqrt(2).

			// Change from original implementation: code no longer checks if luminance is > 0.
			// Decisions on luminance are passed to the Frame implementation.

			// test against the z-buffer.  larger 1/z means the pixel is
			// closer to the viewer than what's already plotted.
			if ooz > (*zbuffer)[xp][yp] {
				(*zbuffer)[xp][yp] = ooz
				// Change from original implementation: no longer calculating luminanceIndex.
				frame.Set(xp, yp, L)
			}
		}
	}
}
