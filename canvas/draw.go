package canvas

import (
	"github.com/platinasystems/weeb/r2"
)

func (c *Context) RoundedRect(x, s r2.X, r float64) {
	dx, dy := r2.XY(r, 0), r2.XY(0, r)
	dz := s / 2
	z := x + dz
	dzc := dz.Conj()

	// Upper left, upper right, lower left, lower right corners.
	ul, ur, ll, lr := z-dz, z+dzc, z-dzc, z+dz

	c.BeginPath()
	c.MoveTo(ul + dx)
	c.ArcTo(ur, ur+dy, r)
	c.ArcTo(lr, lr-dx, r)
	c.ArcTo(ll, ll-dy, r)
	c.ArcTo(ul, ul+dx, r)
}
