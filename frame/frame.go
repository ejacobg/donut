package frame

import "math"

type Frame interface {
	// Apply a luminance value to a point in the frame.
	// If the luminance value was applied, function returns true.
	Set(x, y int, L float64) bool
}

type SetFunc[T any] func(L float64) (T, bool)

// Properties of the torus to be rendered.
type Torus struct {
	// Controls the resolution of the cross-sectional circle.
	ThetaSpacing float64

	// Controls the resolution of the revolution of the circle.
	PhiSpacing float64

	// Radius of the cross-sectional circle.
	R1 float64

	// Distance from axis of revolution (Y) to the center of the cross-sectional circle.
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
	// Obsolete, TX = R2, TY = 0 (always, or else the revolution is messed up), TZ = K2
	// TX, TY, TZ float64

	// Distance from camera to torus. Essentially equal to TZ since camera is at origin.
	K2 float64

	// Light source vector. Technially should be normalized, but up to the implementation.
	LX, LY, LZ float64
}

// Optional. Automatically sets K1 such that the given torus stays within the screen.
func (s *Scene) CalculateK1(R1, R2 float64) {
	min := s.Width
	if s.Height < min {
		min = s.Height
	}
	s.K1 = float64(min) * s.K2 * 3 / (8 * (R1 + R2))
}

func newZBuffer(width, height int) *[][]float64 {
	matrix := make([][]float64, height)
	for i := range matrix {
		matrix[i] = make([]float64, width)
	}
	return &matrix
}

// Renders the scene onto the frame. Assumes that the frame is sized correctly.
func Render(frame Frame, t Torus, s Scene) {
	cosA, sinA := math.Cos(t.A), math.Sin(t.A)
	cosB, sinB := math.Cos(t.B), math.Sin(t.B)

	zbuffer := newZBuffer(s.Width, s.Height)

	// theta goes around the cross-sectional circle of a torus
	for theta := 0.0; theta < 2.0*math.Pi; theta += t.ThetaSpacing {
		cosTheta, sinTheta := math.Cos(theta), math.Sin(theta)

		// phi goes around the center of revolution of a torus
		for phi := 0.0; phi < 2.0*math.Pi; phi += t.PhiSpacing {
			cosPhi, sinPhi := math.Cos(phi), math.Sin(phi)

			// the x,y coordinate of the circle, before revolving (factored out of the above equations)
			circleX := t.R2 + t.R1*cosTheta
			circleY := t.R1 * sinTheta

			// final 3D (x,y,z) coordinate after rotations, directly from our math above
			x := circleX*(cosB*cosPhi+sinA*sinB*sinPhi) - circleY*cosA*sinB
			y := circleX*(sinB*cosPhi-sinA*cosB*sinPhi) + circleY*cosA*cosB
			z := s.K2 + cosA*circleX*sinPhi + circleY*sinA
			ooz := 1 / z // "one over z"

			// x and y projection.  note that y is negated here, because y
			// goes up in 3D space but down on 2D displays.
			xp := int(float64(s.Width)/2 + s.K1*ooz*x)
			yp := int(float64(s.Height)/2 - s.K1*ooz*y)

			// Change from original implementation: using the general form of the luminance dot product.
			L := s.LX*(cosB*cosTheta*cosPhi-sinB*(cosA*sinTheta-sinA*cosTheta*sinPhi)) +
				s.LY*(cosPhi*cosTheta*sinB+cosB*(cosA*sinTheta-cosTheta*sinA*sinPhi)) +
				s.LZ*(cosA*cosTheta*sinPhi+sinA*sinTheta)

			// Change from original implementation: code no longer checks if luminance is > 0.
			// Decisions on luminance are passed to the Frame implementation.

			// test against the z-buffer.  larger 1/z means the pixel is
			// closer to the viewer than what's already plotted.
			if ooz > (*zbuffer)[xp][yp] {
				// Change from original implementation: no longer calculating luminanceIndex.
				if frame.Set(xp, yp, L) {
					// Only update the buffer if something was plotted.
					(*zbuffer)[xp][yp] = ooz
				}
			}
		}
	}
}
