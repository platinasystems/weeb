package canvas

import (
	"github.com/platinasystems/weeb/r2"
)

// Draw rounded rectangle centered about given point.
func (c *Context) CenteredRoundedRect(center, size r2.X, cornerRadius float64) {
	r := cornerRadius
	z := center

	dx, dy := r2.XY(r, 0), r2.XY(0, r)
	dz := size / 2
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

func (c *Context) RoundedRect(x, size r2.X, cornerRadius float64) {
	c.CenteredRoundedRect(x+size/2, size, cornerRadius)
}
