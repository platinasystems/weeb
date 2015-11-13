// The plane: R^2
package r2

import (
	"math"
)

// (x,y) coordinate as complex number for easy arithmetic.
type X complex128

// Angles in units of 2pi/2^31 using modular arithmetic.
type Angle uint32

const AngleMax = 1 << 31

func (a Angle) Radians() float64 {
	n := float64(a&(AngleMax-1)) / float64(AngleMax/2)
	if a == AngleMax {
		n = 2
	}
	return n * math.Pi
}

func (x X) X() float64  { return real(x) }
func (x X) Y() float64  { return imag(x) }
func XY(x, y float64) X { return X(complex(x, y)) }
func XY1(x float64) X   { return XY(x, x) }
func Conj(x X) X        { return XY(real(x), -imag(x)) }
func Norm(x X) float64  { return math.Hypot(x.X(), x.Y()) }

type Rect struct{ X, Size X }
